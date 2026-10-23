# smips

smips는 단순하지만 합리적인 구조를 가지진 않은 가상 머신입니다.

각 호출 프레임이 독립적인 레지스터 묶음을 가지며, 부모 프레임의 인자를 자식 프레임으로 전달하는 레지스터 기반 VM입니다. Go API로 명령어를 구성하고 실행합니다.

수정 전 구현에서 설계를 역도출한 과정은 [초기 설계 분석 기록](docs/architecture.md)에 남겨두었습니다. 현재 동작과 호출 규약은 이 README와 실행 코드를 기준으로 합니다.

## Op Code

smips는 두 가지 형태의 opcode를 가집니다.

`8 bit data, 16 bit ra, 16 bit rb, 16 bit rd, 8 bit op`

`40 bit wide-data, 16 bit rd, 8 bit op`

첫번째 케이스는 8 비트 데이터, 16 비트 레지스터 주소 3개, 8 비트 오퍼레이션가집니다. 이 케이스의 opcode는 거의 대부분의 연산에 쓰입니다. 

두번째 케이스는 40 비트 데이터와 16 비트 레지스터 주소 1개, 8 비트 오퍼레이션을 가지며, 이 경우는 제한된 연산에만 쓰입니다.

`NewWideOpCode`는 40비트를 넘는 즉시값을 거부합니다. 메모리·파일 입출력 명령의 `data`는 0~8바이트의 크기이며, 0은 아무 바이트도 전송하지 않습니다. `OpEq`, `OpNeq`는 조건이 참일 때 `rd`가 가리키는 레지스터의 값으로 분기합니다. 점프 대상은 0부터 시작하는 명령어 인덱스입니다.

각 연산에 대한 코드는 [op.go](opcode/op.go) 파일에 주석과 함께 작성되어 있습니다.

## Memory

smips의 메모리는 고 언어의 해시맵을 그대로 사용하여 오프셋과 함께 나타냅니다.

메모리 최소 단위는 `uint8`이고 주소값과 내부 오프셋으로 접근할 수 있습니다.

자세한 사항은 [memory.go](memory/memory.go) 파일에 주석과 함께 작성되어 있습니다.

## Register

smips의 레지스터는 극단적으로 단순한 형태를 취합니다. 

프레임마다 총 65,536개의 `uint64` 저장 공간을 가집니다. 예약 주소 0, 1, 2, 3의 초기값은 각각 0, 1, 2, 4이며 쓰기도 가능합니다. 나머지는 각각 16,383개인 네 영역으로 나뉩니다.

| 영역 | 명령어 주소 | 주소 상수 |
| --- | --- | --- |
| params | 4~16386 | `register.ParamsBase` |
| temporary | 16387~32769 | `register.TemporaryBase` |
| extra | 32770~49152 | `register.ExtraBase` |
| returns | 49153~65535 | `register.ReturnsBase` |

`Hi`는 마지막 temporary 레지스터(`register.ExtraBase - 1`)와 같은 저장 공간입니다. 일반 레지스터 접근은 항상 최상단 프레임을 대상으로 합니다.

호출 전 부모는 `extra`에 인자를 준비합니다. `Stack.Call`은 이를 새 프레임의 `params`로 복사합니다. `Stack.Plunder`는 자식의 `returns`를 부모의 `extra`로 복사하고 자식 프레임을 풀에 반환합니다. 루트 프레임은 `Plunder`로 제거할 수 없습니다.

자세한 사항은 [reg.go](register/reg.go) 파일에 주석과 함께 작성되어 있습니다.

## Process

smips의 실행 단위입니다.

`Process`는 명령어 배열, PC, 레지스터 프레임 스택, 메모리, 파일 테이블을 소유합니다. `Eval()`은 한 명령어, `Execute()`는 명령어 배열 끝까지 실행합니다. 실행 후에는 루트 프레임 하나가 남습니다.

`AddArgs`는 정수·핸들 등의 원시 값을, `AddStringArgs`는 문자열을 메모리에 저장한 핸들을 순서대로 루트 `params`에 추가합니다. 문자열 표현은 빅엔디언 8바이트 rune 개수와 rune당 4바이트입니다. 실행 결과는 `GetReturn(index)`로 읽습니다.

`AddLib`에 등록한 `Process`는 내부 호출용 코드 소스입니다. `OpLib`의 `ra`는 등록 목록의 인덱스이며, 호출자의 메모리·파일 테이블을 유지한 채 새 레지스터 프레임에서 등록된 코드를 실행합니다. 중첩 호출은 해당 코드 소스의 라이브러리 목록을 사용합니다. 호출할 때마다 PC는 0부터 시작하며 코드 소스의 기존 인자·레지스터·메모리·파일 상태는 사용하거나 변경하지 않습니다. 외부 코드 로더는 제공하지 않습니다.

정상 복귀하면 부모의 PC와 프레임이 복원됩니다. 자식 실행 중 패닉이 나도 프레임과 실행 위치는 복원되지만 자식의 미완성 반환값은 전달하지 않습니다. 메모리 변경이나 파일 쓰기는 되돌리지 않습니다.

자세한 사항은 [process.go](process/process.go) 파일에 주석과 함께 작성되어 있습니다.

## 실행 예제

```go
package main

import (
    "fmt"

    "github.com/snowmerak/smips/opcode"
    "github.com/snowmerak/smips/process"
    "github.com/snowmerak/smips/register"
)

func main() {
    child := process.New(opcode.NewRegularOpCode(
        0, register.ParamsBase, 1, register.ReturnsBase, uint8(opcode.OpAdd),
    ))
    p := process.New(
        opcode.NewRegularOpCode(0, register.ParamsBase, 0, register.ExtraBase, uint8(opcode.OpAdd)),
        opcode.NewRegularOpCode(0, 0, 0, 0, uint8(opcode.OpLib)),
        opcode.NewRegularOpCode(0, register.ExtraBase, 0, register.ReturnsBase, uint8(opcode.OpAdd)),
    ).AddArgs(41).AddLib(child)
    p.Execute()
    fmt.Println(p.GetReturn(0)) // 42
}
```

같은 공개 API 사용 경로를 [실행 가능한 예제 테스트](process/process_integration_test.go)로 검증합니다.

## 입출력과 오류 처리

파일 쓰기는 실제 쓴 바이트 수를 반환하며 시스템 호출은 이를 `Hi`에 저장합니다. 읽기는 실제로 읽은 바이트만 빅엔디언 정수로 결합합니다. EOF에서 읽은 데이터가 없으면 0을 반환하므로, 현재 API는 EOF와 값 0을 별도로 구분하지 않습니다. EOF 외의 입출력 오류, 잘못된 opcode·크기·주소 등은 Go 패닉으로 처리합니다.

표준 입출력 핸들 0, 1, 2는 호스트 소유이므로 `CloseFile`이 닫거나 재사용하지 않으며 0을 반환합니다. 일반 파일의 닫기 성공값은 1입니다.

메모리의 범위 밖 읽기는 0, 범위 밖 쓰기는 패닉입니다. `OpStore`는 전체 쓰기 범위를 먼저 검사합니다. 해제 핸들은 재사용되며, 세대 검사나 GC는 없습니다. 파일 접근은 호스트 OS 권한을 사용하고 열린 파일의 정리는 명시적인 `CloseFile` 호출이 필요합니다. 같은 `Process`의 동시 실행이나 실행 시간 제한은 지원하지 않습니다.

레지스터 주소와 디코더는 초기 구현의 불일치를 수정했습니다. 기존 숫자 주소를 사용하는 코드는 위 주소 상수에 맞춰 확인해야 합니다.

## 검증

Go 1.18 이상에서 다음 명령으로 검사합니다.

```console
go test ./...
go vet ./...
go build ./...
```

테스트는 명령어 비트 배치, 전체 레지스터 주소 공간, 중첩·반복 호출과 실패 복귀, 메모리, 문자열, 실제 임시 파일을 사용하는 시스템 호출을 포함합니다. 인코딩·문자열에는 Go fuzz 타깃도 있습니다. race detector는 지원되는 OS·아키텍처에서 `go test -race ./...`로 실행할 수 있습니다.
