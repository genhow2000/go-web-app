<template>
  <div class="feature-item" @click="$emit('click')">
    <div class="feature-icon">{{ feature.icon }}</div>
    <div class="feature-title">{{ feature.title }}</div>
    <div class="feature-desc">{{ feature.description }}</div>
  </div>
</template>

<script>
export default {
  name: 'FeatureItem',
  props: {
    feature: {
      type: Object,
      required: true
    }
  },
  emits: ['click']
}
</script>

<style scoped>
.feature-item {
  text-align: center;
  padding: 20px;
  border-radius: 10px;
  background: #f7fafc;
  transition: background 0.3s ease;
  cursor: pointer;
}

.feature-item:hover {
  background: #edf2f7;
}

.feature-icon {
  font-size: 2rem;
  margin-bottom: 10px;
}

.feature-title {
  font-weight: bold;
  color: #2d3748;
  margin-bottom: 5px;
}

.feature-desc {
  color: #718096;
  font-size: 0.9rem;
}
</style>
