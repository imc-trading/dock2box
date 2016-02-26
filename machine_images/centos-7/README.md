# CentOS Release

Build will default to CentOS 7 rolling release, if you want a specific release specify REL= either 7.0, 7.1 or 7.2.

```bash
make push
```

Or:

```bash
make push REL=<release>
```

## No Cache

If you want to build with no cache add NOCACHE=1.

```bash
make push NOCACHE=1
```
