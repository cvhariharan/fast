# Fast
A command line utility to quickly check download and (upload) speeds. The final aim of this project is to build something like what [fast.com](https://fast.com) provides as outlined in [this](https://netflixtechblog.com/building-fast-com-4857fe0f8adb) blog post.

## Build Instructions
```bash
go install
fast
```
Output
```bash
Current speed: 425 kB/s            
Average speed: 250 kB/s
Started 21 seconds ago
```
## TODO
- [ ] Measure upload speeds
- [ ] Latency measurements
- [ ] Deploy test files over a CDN