### RUN SERVER
```
go run ./cmd/server
```
### RUN TAILWIND

```
 npx @tailwindcss/cli -i ./public/src/input.css -o ./public/src/output.css --watch
 ```

### RUN MODE DEV / MOBILE

```
export DEV_DISABLE_SECURE_COOKIE=1

```

### RUN TAILWIND

#### 1. Install bun dependencies

``` bash
bun install
```

#### 2. Run if any vue file is modified

``` bash
bun run dev
```
