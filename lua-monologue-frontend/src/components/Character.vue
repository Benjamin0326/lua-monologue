<script setup>
defineProps({
    name: String,
    mod: String
});
</script>

<template>
    <div class="character">
        <h2>{{ name }}</h2>
        <p>현재 감정 상태: {{  mod }}</p>
    </div>
</template>

<style>
.character {
    border: 1px solid #ddd;
    padding: 10px;
    border-radius: 10px;
    text-align: center;
}
</style>