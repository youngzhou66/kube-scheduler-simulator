import { V1Deployment } from "@kubernetes/client-node";
import { reactive, inject } from "@nuxtjs/composition-api";
import { DeploymentAPIKey } from "~/api/APIProviderKeys";
import {
  createResourceState,
  addResourceToState,
  modifyResourceInState,
  deleteResourceInState,
} from "./helpers/storeHelper";
import { WatchEventType } from "@/types/resources";

type stateType = {
  selectedDeployment: selectedDeployment | null;
  deployments: V1Deployment[];
  lastResourceVersion: string;
};

export type selectedDeployment = {
  // isNew represents whether this is a new one or not.
  isNew: boolean;
  item: V1Deployment;
  resourceKind: string;
  isDeletable: boolean;
};


export default function deploymentStore() {
  const state: stateType = reactive({
    selectedDeployment: null,
    deployments: [],
    lastResourceVersion: "",
  });
  const deploymentAPI = inject(DeploymentAPIKey);
  if (!deploymentAPI) {
    throw new Error(`${DeploymentAPIKey.description} is not provided`);
  }
  return {
    get deployments() {
      return state.deployments
    },
    get count(): number {
      return state.deployments.length;
    },

    get selected() {
      return state.selectedDeployment;
    },

    select(d: V1Deployment | null, isNew: boolean) {
      if (d !== null) {
        state.selectedDeployment = {
          isNew: isNew,
          item: d,
          resourceKind: "Deployment",
          isDeletable: true,
        };
      }
    },

    // initList calls list API, and stores current resource data and lastResourceVersion.
    async initList() {
      const list = await deploymentAPI.listDeployment();
      state.deployments = createResourceState<V1Deployment>(list.items)
      state.lastResourceVersion = list.metadata?.resourceVersion!;
    },
    get lastResourceVersion() {
      return state.lastResourceVersion;
    },
    async setLastResourceVersion(pv: V1Deployment) {
      state.lastResourceVersion =
          pv.metadata!.resourceVersion || state.lastResourceVersion;
    },
    // watchEventHandler handles each notified event.
    async watchEventHandler(eventType: WatchEventType, ns: V1Deployment) {
      switch (eventType) {
        case WatchEventType.ADDED: {
          state.deployments = addResourceToState(state.deployments, ns);
          break;
        }
        case WatchEventType.MODIFIED: {
          state.deployments = modifyResourceInState(state.deployments, ns);
          break;
        }
        case WatchEventType.DELETED: {
          state.deployments = deleteResourceInState(state.deployments, ns);
          break;
        }
        default:
          break;
      }
    },
    resetSelected() { },
    async apply(d: V1Deployment) {
      await deploymentAPI.createDeployment(d)
    },
    async fetchSelected() { },
    async delete(d: V1Deployment) {
      if (d.metadata?.name) {
        await deploymentAPI.deleteDeployment(d.metadata.name);
      } else {
        throw new Error(
            "failed to delete deployment: deployment should have metadata.name"
        )
      }
    },


  };
}

export type DeploymentStore = ReturnType<typeof deploymentStore>;