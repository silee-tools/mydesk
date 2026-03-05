# mydesk

macOS config backup & sync tool (Mackup alternative)

## 프로젝트 구조

```
main.go                  # CLI 진입점 (서브커맨드 디스패처)
cmd/                     # 서브커맨드 핸들러 (link, unlink, diff, sync, setup, init)
internal/
  ui/color.go            # ANSI 컬러 출력 (NO_COLOR 지원)
  config/config.go       # 전역 설정, 경로/변수 해석, 설정 레포 탐색
  exec/runner.go         # 외부 명령 실행 (dry-run 지원)
  native/                # 네이티브 디렉토리 규약 (home→~, config→~/.config 등)
    native.go            # 규약 정의
    scanner.go           # 디렉토리 스캔 → LinkEntry 변환
  linker/
    config.go            # links.conf 파싱
    linker.go            # 심볼릭 링크 생성/제거 (백업 포함)
  drift/detector.go      # 드리프트 감지
  provision/             # 프로비저닝 모듈 (brew, vscode, omz, mise, defaults)
```

## 개발

- Language: Go
- Task Runner: mise
- Build: `mise run build`
- Test: `mise run test`
- Lint: `mise run lint`
- Format: `mise run fmt`
- Install: `mise run install`

## 핵심 개념

- **네이티브 디렉토리**: home/, config/, ssh/, vscode/, brew/, macos/, omz/ — 자동 감지
- **links.conf**: 네이티브로 커버되지 않는 추가 심볼릭 링크 매핑
- **설정 레포**: 사용자가 별도 레포에서 관리하는 dotfiles (이 레포가 아님)
- 글로벌 플래그는 서브커맨드 앞에 위치: `mydesk --dry-run link`
