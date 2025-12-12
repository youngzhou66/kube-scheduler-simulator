<template>
  <div>
    <slot />
  </div>
</template>

<script lang="ts">
import { provide, defineComponent } from "@nuxtjs/composition-api";
import DeploymentStoreKey from "~/components/StoreKey/DeploymentStoreKey";
import deploymentStore from "~/store/deployment";

export default defineComponent({
  setup() {
    provide(DeploymentStoreKey, deploymentStore());
    return {};
  },
});
</script>
