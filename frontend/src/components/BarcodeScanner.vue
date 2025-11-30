<template>
  <div class="scanner-container">
    <div id="reader" width="600px"></div>
    <van-button block type="danger" @click="$emit('close')">Close Scanner</van-button>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, defineEmits } from 'vue';
import { Html5QrcodeScanner } from "html5-qrcode";

const emit = defineEmits(['scan', 'close']);
let scanner = null;

onMounted(() => {
  scanner = new Html5QrcodeScanner(
    "reader",
    { fps: 10, qrbox: { width: 250, height: 250 } },
    /* verbose= */ false
  );
  scanner.render(onScanSuccess, onScanFailure);
});

function onScanSuccess(decodedText, decodedResult) {
  // handle the scanned code as you like, for example:
  console.log(`Code matched = ${decodedText}`, decodedResult);
  emit('scan', decodedText);
  if (scanner) {
      scanner.clear();
  }
}

function onScanFailure(error) {
  // handle scan failure, usually better to ignore and keep scanning.
  // for example:
  // console.warn(`Code scan error = ${error}`);
}

onBeforeUnmount(() => {
    if (scanner) {
        scanner.clear().catch(error => {
            console.error("Failed to clear html5-qrcode scanner. ", error);
        });
    }
});
</script>

<style scoped>
.scanner-container {
  padding: 16px;
}
</style>
