<template>
  <div class="tech-category">
    <h3>{{ category.title }}</h3>
    <div v-for="item in category.items" :key="item" class="tech-item">
      {{ item }}
    </div>
  </div>
</template>

<script>
export default {
  name: 'TechCategory',
  props: {
    category: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped>
.tech-category h3 {
  color: #4a5568;
  margin-bottom: 15px;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 5px;
}

.tech-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  color: #718096;
}

.tech-item::before {
  content: "✓";
  color: #48bb78;
  font-weight: bold;
  margin-right: 10px;
}
</style>
