import axiosIns from "@/lib/axios";
import { SupplierRequest, RegisterRequest, User } from "@/lib/types";

export const RegisterCustomer = async (profile: RegisterRequest) => {
  try {
    const response = await axiosIns.post("/api/v1/customer", profile);
    console.log(response.data);
    return response.data.message;
  } catch (error: any) {
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while creating account"
      );
    }
    throw new Error("An error has occured while creating account");
  }
};

export const RegisterSupplier = async(profile: SupplierRequest) => {
  try {
    // Filter out fields that backend doesn't expect
    const requestPayload = {
      tin: profile.tin,
      latitude: profile.latitude,
      longitude: profile.longitude,
      general_zone: profile.general_zone,
      region: profile.region,
      woreda: profile.woreda,
      licence_url: profile.licence_url || "",
      first_name: profile.first_name,
      last_name: profile.last_name,
      email: profile.email,
      phone: profile.phone,
      password: profile.password,
      // username will be set by backend from phone
    };
    const response = await axiosIns.post("/api/v1/supplier", requestPayload);
    console.log(response.data);
    // Backend returns user_id as supplier id
    return response.data.body.supplier;
  } catch (error: any) {
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while creating account"
      );
    }
    throw new Error("An error has occured while creating account");
  }
}

export const InitResetPassword = async (email: string) => {
  try {
    console.log(email);
    const response = await axiosIns.post("/user/init_reset", {
      email: email,
    });
    return response.data.message;
  } catch (error: any) {
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while reseting the password"
      );
    }
    throw new Error("An error has occured while reseting the password");
  }
};

export const ResetPassword = async (token: string, password: string) => {
  try {
    const response = await axiosIns.post(`/auth/reset/${token}`, {
      password: password,
    });
    return response.data.message;
  } catch (error: any) {
    console.log(error);
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while creating the product"
      );
    }
    throw new Error("An error has occured while creating the product");
  }
};

export const UpdateProfile = async (data: Partial<User>) => {
  try {
    const response = await axiosIns.patch("/user", data);
    console.log(response.data);
    return response.data.message;
  } catch (error: any) {
    console.log(error);
    if (error.response) {
      throw new Error(
        error.response.data.message ||
        "An error has occured while updating your profile"
      );
    }
    throw new Error("An error has occured while updating your profile");
  }
};


// Subscription payment is no longer handled client-side.