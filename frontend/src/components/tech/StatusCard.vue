<template>
  <div class="status-card" @click="$emit('click')">
    <h3>
      <span :class="['status-icon', status.icon]"></span>
      {{ status.title }}
      <span v-if="status.isRealtime" class="realtime-badge">即時數據</span>
    </h3>
    <div class="status-value">{{ status.value }}</div>
    <div class="status-detail">{{ status.detail }}</div>
    <div class="click-hint">點擊查看詳細數據</div>
  </div>
</template>

<script>
export default {
  name: 'StatusCard',
  props: {
    status: {
      type: Object,
      required: true
    }
  },
  emits: ['click']
}
</script>

<style scoped>
.status-card {
  background: white;
  border-radius: 15px;
  padding: 25px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
}

.status-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 40px rgba(0,0,0,0.15);
}

.status-card h3 {
  color: #4a5568;
  margin-bottom: 15px;
  font-size: 1.3rem;
  display: flex;
  align-items: center;
}

.status-icon {
  width: 24px;
  height: 24px;
  margin-right: 10px;
  border-radius: 50%;
  display: inline-block;
}

.status-online { background: #48bb78; }
.status-database { background: #4299e1; }
.status-api { background: #ed8936; }
.status-cloud { background: #9f7aea; }

.status-value {
  font-size: 2rem;
  font-weight: bold;
  color: #2d3748;
  margin-bottom: 10px;
}

.status-detail {
  color: #718096;
  font-size: 0.9rem;
}

.realtime-badge {
  display: inline-block;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.7rem;
  font-weight: 600;
  margin-left: 8px;
  animation: pulse 2s infinite;
}

.click-hint {
  font-size: 0.8rem;
  color: #667eea;
  font-weight: 500;
  margin-top: 8px;
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

.status-card:hover .click-hint {
  opacity: 1;
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
  100% {
    opacity: 1;
  }
}
</style>
