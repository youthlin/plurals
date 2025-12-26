# plurals

手写词法分析、语法分析实现 GNU gettext 的表达式解析(`Plural-Forms`)

- 表达式规则 https://www.gnu.org/software/gettext/manual/html_node/Plural-forms.html#index-specifying-plural-form-in-a-PO-file
- 产生式见 `token.go`
- 词法分析 `lex.go`
- 语法树构建 `parse.go`
- 语法树节点定义 `expression.go`
- 参考仓库: https://github.com/ojii/gettext.go, https://github.com/leonelquinteros/gotext
- 用 antlr 实现: https://github.com/youthlin/t

