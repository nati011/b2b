import axiosIns from "@/lib/axios";
import { BuySubscriptionRequest, DistributorRequest, RegisterRequest, User } from "@/lib/types";

export const RegisterRetailer = async (profile: RegisterRequest) => {
  try {
    const response = await axiosIns.post("/api/v1/retailer", profile);
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

export const RegisterDistributor = async(profile: DistributorRequest) => {
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
    const response = await axiosIns.post("/api/v1/distributor", requestPayload);
    console.log(response.data);
    // Backend returns user_id as distributor id
    return response.data.body.distributor;
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


export const BuySubscription = async (profile: BuySubscriptionRequest) => {
  try {
    const response = await axiosIns.post("/api/v1/subscription", profile);
    console.log(response)
    // Backend returns: {body: {id: {id: ..., checkout_url: ..., tx_ref: ...}}}
    // The "id" key contains the SubscribeResponse object
    const subscribeResponse = response.data.body.id;
    return {
      id: subscribeResponse.id,
      checkout_url: subscribeResponse.checkout_url,
      tx_ref: subscribeResponse.tx_ref
    };
  } catch (error: any) {
    if (error.response) {
      console.log(error.response)
      throw new Error(
        error.response.data.message ||
        "An error has occured while buying subscription message"
      );
    }
    throw new Error("An error has occured while buying subscription");
  }
};