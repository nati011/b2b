'use server'
import axiosIns from "@/app/libs/axios";
import { Resource } from "@/app/libs/types";
import { withErrorHandling } from "@/app/libs/error-handling";

const ADD_RESOURCE_COMMAND = "add_resource";
const REMOVE_RESOURCE_COMMAND = "remove_resource";

interface RoleCommandRequest {
  roleId: number;
  resourceId: number;
  command: typeof ADD_RESOURCE_COMMAND | typeof REMOVE_RESOURCE_COMMAND;
}

interface RoleCommandResponse {
  success: boolean;
  message: string;
  data?: {
    resource: number;
  };
}



export const fetchResources = async () => {
    try {
        // Todo: find a way to load on scrollable area
        const response = await axiosIns.get("/resource?limit=100");
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const fetchResourceById = async (id: number) => {
    try {
        const response = await axiosIns.get(`/resource?id=${id}`);
        return response.data;
    } catch (error) {
        throw error;
    }
}

export const createResource = async (resourceData: Partial<Resource>) => {
    try {
        const response = await axiosIns.post("/resource", resourceData);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error occurred while creating the resource";
        }
        throw "An error occurred while creating the resource";
    }
}

export const updateResource = async (resourceData: Partial<Resource>) => {
    try {
        const response = await axiosIns.put("/resource", resourceData);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error occurred while updating the resource";
        }
        throw "An error occurred while updating the resource";
    }
}

export const deleteResource = async (id: number) => {
    try {
        const response = await axiosIns.delete(`/resource?id=${id}`);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error occurred while deleting the resource";
        }
        throw "An error occurred while deleting the resource";
    }
}


export const executeRoleCommand = async ({
    roleId,
    resourceId,
    command
  }: RoleCommandRequest): Promise<RoleCommandResponse> => {
    return withErrorHandling(async () => {
      if (command !== ADD_RESOURCE_COMMAND && command !== REMOVE_RESOURCE_COMMAND) {
        throw new Error("Invalid command. Must be 'add_resource' or 'remove_resource'");
      }

      if (!roleId || roleId <= 0) {
        throw new Error("Invalid role ID");
      }
  
      if (!resourceId || resourceId <= 0) {
        throw new Error("Invalid resource ID");
      }

      const url = `/role/${roleId}?command=${command}&resource_id=${resourceId}`;
  
      try {
        const response = await axiosIns.patch(url);
        console.log(url)
        
        return {
          success: true,
          message: command === ADD_RESOURCE_COMMAND 
            ? "Resource added to role successfully" 
            : "Resource removed from role successfully",
          data: response.data
        };
      } catch (error: any) {
        if (error.response?.status === 404) {
          throw new Error("Role or resource not found");
        }
        
        throw error;
      }
    });
  };

  export const addResourceToRole = async (roleId: number, resourceId: number): Promise<RoleCommandResponse> => {
    return executeRoleCommand({
      roleId,
      resourceId,
      command: ADD_RESOURCE_COMMAND
    });
  };
  
  export const removeResourceFromRole = async (roleId: number, resourceId: number): Promise<RoleCommandResponse> => {
    return executeRoleCommand({
      roleId,
      resourceId,
      command: REMOVE_RESOURCE_COMMAND
    });
  };
  

  export const addMultipleResourcesToRole = async (
    roleId: number, 
    resourceIds: number[]
  ): Promise<RoleCommandResponse[]> => {
    return withErrorHandling(async () => {
      const results = await Promise.allSettled(
        resourceIds.map(resourceId => addResourceToRole(roleId, resourceId))
      );
      
      const responses: RoleCommandResponse[] = [];
      const errors: string[] = [];
      
      results.forEach((result, index) => {
        if (result.status === 'fulfilled') {
          responses.push(result.value);
        } else {
          errors.push(`Resource ${resourceIds[index]}: ${result.reason.message}`);
        }
      });
      
      if (errors.length > 0) {
        throw new Error(`Some operations failed: ${errors.join(', ')}`);
      }
      
      return responses;
    });
  };
  
  export const removeMultipleResourcesFromRole = async (
    roleId: number, 
    resourceIds: number[]
  ): Promise<RoleCommandResponse[]> => {
    return withErrorHandling(async () => {
      const results = await Promise.allSettled(
        resourceIds.map(resourceId => removeResourceFromRole(roleId, resourceId))
      );
      
      const responses: RoleCommandResponse[] = [];
      const errors: string[] = [];
      
      results.forEach((result, index) => {
        if (result.status === 'fulfilled') {
          responses.push(result.value);
        } else {
          errors.push(`Resource ${resourceIds[index]}: ${result.reason.message}`);
        }
      });
      
      if (errors.length > 0) {
        throw new Error(`Some operations failed: ${errors.join(', ')}`);
      }
      
      return responses;
    });
  };