# AI服务配置指南

## 环境变量配置

在部署应用之前，需要设置以下环境变量：

### 必需的AI服务配置

```bash
# AI服务提供商配置
AI_PRIMARY_PROVIDER=groq
AI_FALLBACK_PROVIDER=gemini
AI_SIMULATION_PROVIDER=simulation

# Groq API配置
GROQ_API_KEY=your_groq_api_key_here
GROQ_API_URL=https://api.groq.com/openai/v1/chat/completions
GROQ_MODEL=llama-3.1-8b-instant
GROQ_MAX_TOKENS=100
GROQ_TEMPERATURE=0.7
GROQ_DAILY_LIMIT=10000

# Google Gemini API配置
GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_API_URL=https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent
GEMINI_MODEL=gemini-2.0-flash
GEMINI_MAX_TOKENS=100
GEMINI_TEMPERATURE=0.7
GEMINI_DAILY_LIMIT=1500
```

### 可选配置

```bash
# Hugging Face API配置（可选）
HF_API_TOKEN=your_huggingface_token_here
HF_API_URL=https://api-inference.huggingface.co/models/microsoft/DialoGPT-small
HF_MODEL=microsoft/DialoGPT-small
HF_MAX_TOKENS=100
HF_TEMPERATURE=0.7
HF_DAILY_LIMIT=1000

# AI服务切换配置
AI_SWITCH_THRESHOLD=0.8
AI_REQUEST_TIMEOUT=30
```

## 如何获取API密钥

### 1. Groq API密钥
1. 访问 [Groq Console](https://console.groq.com/)
2. 注册/登录账户
3. 在API Keys页面创建新的API密钥
4. 复制密钥并设置到 `GROQ_API_KEY` 环境变量

### 2. Google Gemini API密钥
1. 访问 [Google AI Studio](https://aistudio.google.com/)
2. 登录Google账户
3. 创建新项目或选择现有项目
4. 在API Keys页面创建新的API密钥
5. 复制密钥并设置到 `GEMINI_API_KEY` 环境变量

## 部署说明

### Docker部署
```bash
# 设置环境变量
export GROQ_API_KEY=your_groq_api_key
export GEMINI_API_KEY=your_gemini_api_key

# 启动应用
docker-compose up -d
```

### 云平台部署
在云平台（如Google Cloud Run、AWS、Azure等）的环境变量设置中添加上述配置。

## AI服务梯队说明

应用使用智能AI服务梯队：

1. **主要服务**: Groq API - 快速响应，适合实时对话
2. **备用服务**: Google Gemini API - 高质量回复，支持多语言
3. **最终降级**: 本地模拟服务 - 确保服务始终可用

当主要服务失败时，系统会自动切换到备用服务，确保聊天功能不中断。

## 测试AI服务

使用内置的测试工具验证AI服务配置：

```bash
# 构建测试镜像
docker build -t go-ai-test .

# 运行AI服务测试
docker run --rm \
  -e GROQ_API_KEY=your_groq_key \
  -e GEMINI_API_KEY=your_gemini_key \
  -e AI_PRIMARY_PROVIDER=groq \
  -e AI_FALLBACK_PROVIDER=gemini \
  go-ai-test ./test_ai
```

## 监控和日志

应用会记录详细的AI服务使用情况：
- 每个服务的使用次数
- 错误统计
- 切换日志
- 响应时间

查看日志：
```bash
docker logs container_name
```