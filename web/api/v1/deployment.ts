import { V1Deployment, V1DeploymentList } from "@kubernetes/client-node";
import { AxiosInstance } from "axios";

export default function deploymentAPI(
  k8sDeploymentInstance: AxiosInstance,
  instance: AxiosInstance
) {
  return {
    listDeployment: async () => {
      try {
        const res = await k8sDeploymentInstance.get<V1DeploymentList>("/deployments", {});
        return res.data;
      } catch (e: any) {
        throw new Error("failed to list");
      }
    },

    createDeployment: async (req: V1Deployment) => {
      try {
        const res = await instance.post<V1Deployment>("/addDeployment", req, {
          headers: { "Content-Type": "application/json" },
        });
        return res.data;
      } catch (e: any) {
        throw new Error(`failed to create deployment: ${e}`);
      }
    },

    deleteDeployment: async (name: string) => {
      try {
        const res = await instance.delete<V1Deployment>(`/deleteDeployment/default/${name}`, {});
        return res.data;
      } catch (e: any) {
        throw new Error(`failed to delete deployment: ${e}`);
      }
    },

  };
}
export type DeploymentAPI = ReturnType<typeof deploymentAPI>;
