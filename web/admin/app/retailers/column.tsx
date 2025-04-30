export type Retailer = {
    id: number;
    name: string;
    tin: string;
    latitude: string;
    longitude: string
    general_zone: string;
    region: string;
    woreda: string;
    user: UserAccount;
}
export const column = [
    {
        title: 'Id',
        dataIndex: 'Id',
        key: 'Id',
    },
    {
        title: 'Name',
        dataIndex: 'Name',
        key: 'Name',
    },
    {
        title: 'Tin',
        dataIndex: 'tin',
        key: 'Tin'
    },
    {
        title: 'Latitude',
        dataIndex: 'latitude',
        key: 'latitude'
    },
    {
        title: 'Is Active',
        dataIndex: 'IsActive',
        key: 'IsActive',
    },
];


