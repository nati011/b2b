import { redirect } from "next/navigation";

export function handleServerActionError(error: any) {
  console.error("Server action error:", error);
  
  if (error.name === "ForbiddenError") {
    redirect("/forbidden");
  }
  
  if (error.name === "UnauthorizedError" || error.name === "TokenExpiredError" || error.name === "TokenExpiredError") {
    redirect("/auth/signin");
  }
  
  if (error.redirectTo) {
    redirect(error.redirectTo);
  }
  
  throw error;
}

export async function withErrorHandling<T>(
  action: () => Promise<T>
): Promise<T> {
  try {
    return await action();
  } catch (error: any) {
    if (error instanceof Error && error.name === "NEXT_REDIRECT") {
      throw error;
    }

    if (
      error.name === "ForbiddenError" ||
      error.name === "UnauthorizedError" ||
      error.name === "TokenExpiredError" ||
      error.redirectTo
    ) {
      handleServerActionError(error);
      throw error; 
    }

    if (error.response) {
      const status = error.response.status;
      const message = error.response.data?.message || error.message;

      switch (status) {
        case 400:
          throw new Error(`Validation Error: ${message}`);
        case 401:
          redirect("/auth/signin");
        case 422:
          throw new Error(`Invalid Data: ${message}`);
        case 429:
          throw new Error("Too many requests. Please try again later.");
        case 500:
          throw new Error("Server error. Please try again later.");
        default:
          console.log(message);
          throw new Error(message);
      }
    }

    if (error.code === 'ECONNABORTED') {
      throw new Error("Request timeout. Please check your connection and try again.");
    }

    if (error.code === 'NETWORK_ERROR') {
      throw new Error("Network error. Please check your connection and try again.");
    }

    throw error;
  }
}