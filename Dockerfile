FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o server ./cmd/

FROM alpine:latest

LABEL authors="zhxycn"

WORKDIR /app

COPY --from=builder /app/server .

ENV PORT="8080"
ENV DEBUG="false"
ENV DATABASE=""
ENV TENCENT_APP_ID=""
ENV TENCENT_SECRET_ID=""
ENV TENCENT_SECRET_KEY=""
ENV XUNFEI_APP_ID=""
ENV XUNFEI_API_KEY=""
ENV XUNFEI_API_SECRET=""
ENV DIFY_ENDPOINT=""
ENV RESUME_API_KEY=""
ENV INTERVIEW_API_KEY=""
ENV SILICONFLOW_TOKEN=""
ENV SILICONFLOW_MODEL="FunAudioLLM/CosyVoice2-0.5B"
ENV SILICONFLOW_VOICE="FunAudioLLM/CosyVoice2-0.5B:alex"

EXPOSE 8080

CMD ["./server"]