<template>
  <v-card v-if="deployments.length !== 0" class="ma-2" outlined>
    <v-card-title class="mb-1"> Deployments </v-card-title>
    <v-container>
      <v-row no-gutters>
        <v-col v-for="(n, i) in deployments" :key="i" tile cols="auto">
          <v-card class="ma-2" outlined @click="onClick(n)">
            <v-card-title>
              <img
                  src="/node.svg"
                  height="40"
                  alt="d.metadata.name"
                  class="mr-2"
              />
              {{ n.metadata.name }}
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-card>
</template>

<script lang="ts">
import { computed, inject, defineComponent } from "@nuxtjs/composition-api";
import { V1Deployment, V1Node } from "@kubernetes/client-node";
import {} from "../../lib/util";
import DeploymentStoreKey from "~/components/StoreKey/DeploymentStoreKey";

export default defineComponent({
  setup() {
    const dstore = inject(DeploymentStoreKey);
    if (!dstore) {
      throw new Error(`${DeploymentStoreKey.description} is not provided`);
    }

    const deployments = computed(() => dstore.deployments);
    const onClick = (deployment: V1Deployment) => {
      dstore.select(deployment, false);
    };

    return {
      deployments,
      onClick,
    };
  },
});
</script>