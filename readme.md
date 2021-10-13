# gitagger

## create a new tag/release
```
export CI_API_V4_URL=https://gitlab.domain.tld
export CI_BUILD_TOKEN=a...secret...z

#        <project_id>:<git_ref>:<tag>:<message>

gitagger createTag 426:master:v99.9.10:release

```

## files

```
#        <project_id>:<git_ref>:<repo_path>:<content>

gitagger createFile 426:master:etc/example.dot:hello-world

# write a file
gitagger createFile '426:master:etc/example.dot:<@path/file.ext'

# update file from local file
gitagger updateFile '426:master:etc/example.dot:<@path/file.ext'
```
