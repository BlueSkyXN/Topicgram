# Ad Filter / 广告过滤器

## 概述 / Overview

广告过滤器是一个基于文本识别的消息拦截系统，可以自动检测和阻止垃圾广告消息。

The Ad Filter is a text recognition-based message interception system that automatically detects and blocks spam/advertising messages.

## 功能特性 / Features

- **关键词匹配** / Keyword Matching: 支持配置关键词黑名单，自动拦截包含这些关键词的消息
- **通配符关键词** / Wildcard Keywords: 支持使用 `*` 和 `?` 通配符创建灵活的匹配规则 ([详细文档](WildcardFiltering.md))
- **正则表达式** / Regular Expressions: 支持使用正则表达式定义复杂的过滤规则
- **URL 限制** / URL Limits: 可以限制消息中的链接数量，防止广告链接泛滥
- **实时日志** / Real-time Logging: 记录所有被拦截的消息，方便审计和调整规则

## 配置说明 / Configuration

在 `config.json` 的 `Bot` 配置中添加 `AdFilter` 字段：

Add the `AdFilter` field to the `Bot` configuration in `config.json`:

```json
{
    "Bot": {
        "Token": "your_bot_token",
        "GroupId": 123456789,
        "LanguageCode": "zh-hans",
        "WebHook": {
            "Host": "your.domain.com"
        },
        "AdFilter": {
            "enabled": true,
            "keywords": [
                "免费领取",
                "点击领取",
                "限时优惠",
                "立即购买",
                "加微信",
                "添加微信"
            ],
            "wildcard_keywords": [
                "*微信*",
                "*vx*",
                "加*好友",
                "*免费*"
            ],
            "patterns": [
                "a{6,}",
                "\\d{10,}"
            ],
            "url_patterns": [
                "(?i)(https?://|www\\.|t\\.me/)[^\\s]+"
            ],
            "max_url_count": 3
        }
    }
}
```

### 配置项说明 / Configuration Options

- `enabled` (bool): 是否启用广告过滤器 / Whether to enable the ad filter
- `keywords` (array): 关键词黑名单列表（精确匹配）/ List of blocked keywords (exact match)
- `wildcard_keywords` (array): 通配符关键词列表（支持 * 和 ?）/ List of wildcard keywords (supports * and ?) - [详细说明](WildcardFiltering.md)
- `patterns` (array): 正则表达式规则列表 / List of regex patterns
- `url_patterns` (array): URL 匹配规则 / URL matching patterns
- `max_url_count` (int): 允许的最大链接数量 / Maximum allowed URL count

## 默认规则 / Default Rules

如果不配置 `AdFilter`，过滤器默认是禁用的。可以使用以下默认规则启用：

If `AdFilter` is not configured, the filter is disabled by default. You can enable it with these default rules:

### 默认关键词 / Default Keywords

- 免费领取 (Free claim)
- 点击领取 (Click to claim)
- 限时优惠 (Limited time offer)
- 立即购买 (Buy now)
- 加微信 (Add WeChat)
- 添加微信 (Add WeChat)
- 扫码添加 (Scan to add)
- 代理加盟 (Agent/Franchise)
- 兼职 (Part-time job)
- 刷单 (Fake orders)
- 投资理财 (Investment)
- 贷款 (Loan)

### 默认模式 / Default Patterns

- 连续重复字符检测 / Repeated character detection
- 过长数字串检测 / Long number string detection
- URL 数量限制（默认最多 3 个）/ URL count limit (default max 3)

## 工作原理 / How It Works

1. **消息拦截点** / Message Interception Points
   - 用户发送消息到 Bot 时 / When users send messages to the Bot
   - 管理员在 Topic 中回复时 / When admins reply in Topics

2. **检测流程** / Detection Flow
   - 检查消息文本和图片说明 / Check message text and captions
   - 按顺序匹配关键词、正则表达式和 URL 规则 / Match keywords, patterns, and URL rules in order
   - 如果匹配任何规则，消息将被阻止 / If any rule matches, the message is blocked

3. **处理方式** / Handling
   - 删除违规消息 / Delete violating messages
   - 向发送者发送通知 / Send notification to sender
   - 记录日志以供审计 / Log for auditing

## 注意事项 / Notes

1. **误拦截** / False Positives: 请仔细配置关键词和规则，避免误拦截正常消息
2. **性能** / Performance: 规则过多可能影响消息处理性能，建议保持规则精简
3. **管理员豁免** / Admin Bypass: 以 `/` 开头的命令消息不会被过滤
4. **实时调整** / Real-time Adjustment: 修改配置后需要重启 Bot 才能生效

## 示例场景 / Example Scenarios

### 场景 1: 拦截微信广告 / Block WeChat Ads
```json
{
    "keywords": ["加微信", "添加微信", "微信号", "wx"]
}
```

### 场景 2: 拦截多链接消息 / Block Multi-Link Messages
```json
{
    "max_url_count": 1
}
```

### 场景 3: 拦截重复字符 / Block Repeated Characters
```json
{
    "patterns": ["a{5,}", "b{5,}", "c{5,}"]
}
```

## 故障排除 / Troubleshooting

### 过滤器不工作 / Filter Not Working
- 检查 `enabled` 是否设置为 `true`
- 查看日志确认过滤器是否已初始化
- 确认规则语法正确

### 误拦截正常消息 / False Positives
- 检查关键词是否过于宽泛
- 调整正则表达式规则
- 增加 `max_url_count` 的值

### 查看过滤日志 / View Filter Logs
过滤器会在日志中记录所有拦截的消息，查找包含 `[AdFilter]` 的日志行。

The filter logs all intercepted messages. Look for log lines containing `[AdFilter]`.

## 贡献 / Contributing

欢迎提交 Issue 和 Pull Request 来改进广告过滤功能！

Feel free to submit Issues and Pull Requests to improve the ad filtering feature!
