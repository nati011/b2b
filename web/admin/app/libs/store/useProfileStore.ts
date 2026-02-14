import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Profile } from '@/app/libs/types';

interface ProfilesStore {
    profile: Profile;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchProfile: (url?: string) => Promise<void>;
    createProfiles: (ProfilesData: Partial<Profile>) => Promise<void>;
}

const useProfilesStore = create<ProfilesStore>((set) => ({
    profile: {
        id: 0,
        first_name: '',
        last_name: '',
        email: '',
        phone: '',
        username: '',
        dob: ''
    },
    loading: false,
    error: null,
    next: null,
    previous: null,

    fetchProfile: async () => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get('/api/user');
            set({
                profile: response.data.body.Profiles,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch profile', loading: false });
        }
    },

    createProfiles: async (ProfilesData: Partial<Profile>) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post(`/api/user/`, ProfilesData);
            set(state => ({
                profile: response.data.detail,
                loading: false
            }));
        } catch (error) {
            set({ error: 'Failed to create profile', loading: false });
        }
    },

}));

export default useProfilesStore;