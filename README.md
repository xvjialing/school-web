# README

## 根据特定数据表快速生成代码
`bee generate appcode -tables="tablename1,tablename2" -driver=mysql -conn="root:password@tcp(127.0.0.1:3306)/test" -level=3`，其中`level`为生成代码的级别

## 启动时添加自动文档参数：
`bee run -gendoc=true -downdoc=true`
