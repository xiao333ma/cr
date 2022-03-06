
通过 `cr` 来快速打开网页发起 `CR`

目前只支持腾讯工蜂

`cr` 打开 cr 界面
`cr -f` 发起一个到 feature 分支的 cr
`cr -d` 发起一个到 develop 分支的 cr
`cr -s source_branch -t target_branch` 发起一个 source_branch 到 target_branch 的 cr

`cr -p` 指定子目录下的 git 工程的 cr，配合 -f -d 等命令使用，省去目录切换工作