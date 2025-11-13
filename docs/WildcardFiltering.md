# Wildcard Keyword Filtering Guide / 通配符关键词过滤指南

## 概述 / Overview

通配符关键词过滤功能允许使用更灵活的模式匹配来检测广告和垃圾消息。与精确关键词匹配不同，通配符支持使用 `*` 和 `?` 来创建更强大的过滤规则。

Wildcard keyword filtering allows more flexible pattern matching to detect ads and spam messages. Unlike exact keyword matching, wildcards support using `*` and `?` to create more powerful filtering rules.

## 通配符语法 / Wildcard Syntax

### `*` - 匹配任意字符序列 / Match Any Sequence

`*` 通配符可以匹配零个或多个任意字符。

The `*` wildcard matches zero or more characters of any kind.

**示例 / Examples:**

```json
{
    "wildcard_keywords": [
        "*微信*",      // 匹配包含"微信"的任何文本
        "加*好友",     // 匹配"加"后跟任何内容再跟"好友"
        "免费*",       // 匹配以"免费"开头的内容
        "*优惠"        // 匹配以"优惠"结尾的内容
    ]
}
```

**匹配示例 / Match Examples:**

- `*微信*` 匹配:
  - ✅ "请加我的微信"
  - ✅ "微信号12345"
  - ✅ "联系微信"
  - ❌ "请联系我" (不包含"微信")

- `加*好友` 匹配:
  - ✅ "加我好友"
  - ✅ "加个好友吧"
  - ✅ "加她为好友"
  - ❌ "好友申请" (不以"加"开头)

### `?` - 匹配单个字符 / Match Single Character

`?` 通配符匹配恰好一个任意字符。

The `?` wildcard matches exactly one character of any kind.

**示例 / Examples:**

```json
{
    "wildcard_keywords": [
        "v?",          // 匹配v后跟任意单个字符
        "微信?",       // 匹配"微信"后跟任意单个字符
        "vx???"       // 匹配"vx"后跟恰好3个字符
    ]
}
```

**匹配示例 / Match Examples:**

- `v?` 匹配:
  - ✅ "vx号码"
  - ✅ "vv123"
  - ❌ "v" (缺少后续字符)
  - ❌ "vxyz" (多于一个字符)

- `微信?` 匹配:
  - ✅ "微信号"
  - ✅ "微信啊"
  - ❌ "微信" (缺少后续字符)

## 配置示例 / Configuration Examples

### 基础配置 / Basic Configuration

```json
{
    "Bot": {
        "AdFilter": {
            "enabled": true,
            "wildcard_keywords": [
                "*微信*",
                "*vx*",
                "加*好友"
            ]
        }
    }
}
```

### 完整配置 / Full Configuration

```json
{
    "Bot": {
        "AdFilter": {
            "enabled": true,
            "keywords": [
                "免费领取",
                "加微信"
            ],
            "wildcard_keywords": [
                "*微信*",
                "*vx*",
                "加*好友",
                "*免费*",
                "*优惠*码",
                "代理*",
                "*赚钱*"
            ],
            "patterns": [
                "\\d{10,}"
            ],
            "max_url_count": 3
        }
    }
}
```

## 实际应用场景 / Practical Use Cases

### 场景 1: 拦截各种形式的微信广告 / Block Various WeChat Ads

```json
{
    "wildcard_keywords": [
        "*微信*",      // 匹配包含"微信"的任何内容
        "*vx*",       // 匹配WeChat的常见缩写
        "*wx*",       // 另一种常见缩写
        "加*微*",     // 匹配"加"和"微"之间有任何字符
        "*扫码*微*"   // 匹配扫码加微信相关
    ]
}
```

**拦截的消息示例:**
- "加我微信领红包"
- "vx: 123456789"
- "wx号联系我"
- "加个微信好友"
- "扫码加微信"

### 场景 2: 拦截促销和优惠广告 / Block Promotional Ads

```json
{
    "wildcard_keywords": [
        "*免费*",
        "*优惠*",
        "*折扣*",
        "*限时*",
        "*抢购*",
        "立即*"
    ]
}
```

**拦截的消息示例:**
- "免费领取大礼包"
- "限时优惠，立即抢购"
- "8折优惠券"
- "立即购买"

### 场景 3: 拦截招聘和兼职广告 / Block Recruitment Ads

```json
{
    "wildcard_keywords": [
        "*招聘*",
        "*兼职*",
        "*代理*",
        "*加盟*",
        "日赚*",
        "*刷单*"
    ]
}
```

### 场景 4: 组合使用精确和通配符关键词 / Combine Exact and Wildcard Keywords

```json
{
    "keywords": [
        "贷款",
        "投资理财",
        "刷单"
    ],
    "wildcard_keywords": [
        "*贷款*",
        "*投资*",
        "*理财*",
        "*高收益*",
        "日入*元"
    ]
}
```

## 高级技巧 / Advanced Tips

### 1. 使用多个通配符 / Using Multiple Wildcards

```json
{
    "wildcard_keywords": [
        "*微*信*",     // 匹配包含"微"和"信"的文本
        "*免*费*",     // 匹配包含"免"和"费"的文本
        "*?信?"       // 使用?进行更精确的匹配
    ]
}
```

### 2. 前缀和后缀匹配 / Prefix and Suffix Matching

```json
{
    "wildcard_keywords": [
        "代理*",       // 匹配以"代理"开头的所有内容
        "*号码",       // 匹配以"号码"结尾的所有内容
        "vx*"         // 匹配以"vx"开头
    ]
}
```

### 3. 精确长度匹配 / Exact Length Matching

使用多个 `?` 来匹配特定长度的内容。

Use multiple `?` to match specific length content.

```json
{
    "wildcard_keywords": [
        "vx????",      // 匹配"vx"后跟恰好4个字符
        "微信????"     // 匹配"微信"后跟恰好4个字符
    ]
}
```

## 性能注意事项 / Performance Considerations

1. **通配符性能 / Wildcard Performance**: 通配符匹配比精确匹配稍慢，但对于大多数用例来说影响可以忽略不计。

2. **规则数量 / Number of Rules**: 建议保持通配符关键词在 50 个以内以获得最佳性能。

3. **复杂模式 / Complex Patterns**: 对于非常复杂的模式，考虑使用正则表达式 `patterns` 字段。

## 大小写敏感性 / Case Sensitivity

所有通配符匹配都是**不区分大小写**的。

All wildcard matching is **case-insensitive**.

```json
{
    "wildcard_keywords": ["*WeChat*"]
}
```

将匹配 / Will match:
- "WeChat" ✅
- "wechat" ✅
- "WECHAT" ✅
- "WeChaT" ✅

## 与其他过滤器的结合 / Combining with Other Filters

通配符关键词可以与其他过滤器类型结合使用：

Wildcard keywords can be combined with other filter types:

```json
{
    "AdFilter": {
        "enabled": true,
        "keywords": ["spam", "广告"],           // 精确匹配
        "wildcard_keywords": ["*微信*", "vx*"], // 通配符匹配
        "patterns": ["\\d{10,}"],              // 正则表达式
        "max_url_count": 3                      // URL数量限制
    }
}
```

消息将被拦截如果满足**任何一个**条件。

A message will be blocked if it matches **any** of the conditions.

## 测试通配符规则 / Testing Wildcard Rules

在生产环境部署前，建议先测试通配符规则：

Before deploying to production, it's recommended to test wildcard rules:

1. 启用日志记录以查看被拦截的消息
2. 从少量规则开始
3. 监控误报率
4. 根据需要调整规则

## 常见问题 / FAQ

### Q: 通配符和正则表达式有什么区别？

A: 通配符更简单易用，适合基本的模式匹配。正则表达式更强大但也更复杂。通配符使用 `*` 和 `?`，而正则表达式支持更复杂的语法。

### Q: 可以在一个模式中混合使用 * 和 ? 吗？

A: 可以！例如：`"v?*"` 匹配"v"后跟一个字符，然后是任意内容。

### Q: 通配符关键词区分大小写吗？

A: 不区分。所有匹配都是不区分大小写的。

### Q: 如果我想匹配字面的 * 或 ? 字符怎么办？

A: 对于字面字符匹配，使用 `keywords` 字段进行精确匹配，或使用 `patterns` 字段的正则表达式。

## 示例：完整的垃圾消息过滤配置 / Example: Complete Spam Filter Configuration

```json
{
    "Bot": {
        "AdFilter": {
            "enabled": true,
            "keywords": [
                "贷款",
                "理财",
                "刷单"
            ],
            "wildcard_keywords": [
                "*微信*",
                "*vx*",
                "*wx*",
                "加*好友",
                "*免费*",
                "*优惠*",
                "*折扣*",
                "*限时*",
                "代理*",
                "*招聘*",
                "*兼职*",
                "*赚钱*",
                "日赚*",
                "*投资*"
            ],
            "patterns": [
                "\\d{10,}",
                "a{6,}"
            ],
            "url_patterns": [
                "(?i)(https?://|www\\.|t\\.me/)[^\\s]+"
            ],
            "max_url_count": 3
        }
    }
}
```

这个配置将拦截大多数常见的中文垃圾消息和广告。

This configuration will block most common Chinese spam messages and advertisements.
