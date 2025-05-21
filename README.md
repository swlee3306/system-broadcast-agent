# Zeroconf Multi-IP Agent Discovery

이 프로젝트는 동일 네트워크 상의 여러 Agent 간에 서로의 정보를 탐지할 수 있도록 설계된 **Zeroconf(mDNS/Bonjour) 기반 서비스**입니다.  
각 Agent는 자신이 보유한 **모든 사용 가능한 IP 주소와 호스트명**을 네트워크에 광고하며,  
어느 한 Agent가 `/discovery` API를 호출하면 주변 에이전트들을 브라우징하여 정보를 수집합니다.

---

## 🧩 기능 요약

- Zeroconf(mDNS)를 통해 Agent 서비스 자동 탐색
- 각 Agent는 여러 NIC 환경에서 **모든 IPv4 주소를 자동 등록**
- `/discovery` API를 통해 주변 Agent의 IP/호스트명 수집
- Zeroconf TXT 레코드 갱신을 위해 주기적 재등록 수행
- 타임아웃은 API 요청 시 쿼리로 조절 가능 (`?timeout=5`)

---

## 🛠️ 구성 파일

| 파일명         | 설명 |
|----------------|------|
| `main.go`      | Zeroconf 등록 루프 및 HTTP 서버 구동 |
| `zeroconf.go`  | Zeroconf 등록 함수 (멀티 IP 대응) |
| `discovery.go` | `/discovery` API 처리 및 브라우징 |
| `utils.go`     | 모든 사용 가능한 IPv4 수집 로직 |

---

## 🚀 실행 방법

### 1. 의존성 설치
```bash
go mod tidy