the build process:
-

* generate embedded assets => go run generateAssets.go
* build the executable
  * before the build; need to rename the generateAssets.go to generateAssets.go.disable (to exclude it from the build chain)
  * build => go build -o {{target_app_name}}
  * rename the generateAssets.go.disable -> generateAssests.go (revert the 1st process)

---


