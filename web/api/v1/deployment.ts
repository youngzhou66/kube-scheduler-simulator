import { V1Deployment } from "@kubernetes/client-node";
import { AxiosInstance } from "axios";

export default function deploymentAPI(
  instance: AxiosInstance
) {
  return {
    listDeployment: async () => {
      try {
        return {};
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
    }

  };
}
export type DeploymentAPI = ReturnType<typeof deploymentAPI>;
