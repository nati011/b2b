'use client'
import Heading from "@/app/components/breadcrumb";
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { Distributor, DistributorRequest } from '@/app/libs/types';
import dynamic from "next/dynamic";
import 'leaflet/dist/leaflet.css'
import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useParams } from "next/navigation";

const Map = dynamic(
    () => import('@/app/components/map'),
    { ssr: false }
)


export default function DistributorsForm() {
    const {
        distributor,
        distributorUser,
        loading,
        error,
        fetchDistributorUser,
        fetchDistributorDetail
    } = useDistributorsStore()
    const [formData, setFormData] = useState<DistributorRequest>({
        name: "",
        tin: "",
        latitude: "",
        longitude: "",
        general_zone: "",
        region: "",
        woreda: "",
        first_name: "",
        last_name: "",
        email: "",
        phone: "",
        username: "",
        dob: "",
        external_id: "",

    })
    const [markerPosition, setMarkerPosition] = useState<[number, number]>([8.9934609, 38.7714897])
    const pages = [
        {
            "title": "Distributor",
            "href": "/distributor"
        },
        {
            "title": "Detail",
            "href": "/distributor/form"
        }
    ]
    const routeParam = useParams<{ id: string }>();
    useEffect(() => {
        fetchDistributorDetail(parseInt(routeParam.id))
        fetchDistributorUser(distributor.user[0])
    }, [])
    useEffect(() => {
        console.log([distributor])
        console.log([distributor.latitude, distributor.longitude])

        fetchDistributorUser(distributor.user[0])
        setMarkerPosition([parseFloat(distributor.latitude), parseFloat(distributor.longitude)])
    }, [distributor])


    function handleSubmit(): void {
        throw new Error("Function not implemented.");
    }

    return (
        <div className="grid grid-cols-1 gap-4">
            <Heading page={pages} heading="Distributors Details" subheading={`Detail information for ${formData.first_name} ${formData.last_name}`} />
            <Tabs value="business" className="w-full">
                <TabsList className="w-full">
                    <TabsTrigger value="business">Business Information</TabsTrigger>
                    <TabsTrigger value="profile">Profile Information</TabsTrigger>
                </TabsList>
                <TabsContent value="profile">

                    <Card className="rounded-sm border-2 border-gray-200 shadow-none">
                        <CardContent className="">
                            <div className="mb-4">
                                <h3 className="text-xl font-bold">
                                    Profile Information
                                </h3>
                                <p className="text-sm text-muted-foreground">
                                    Profile information for user {distributorUser.first_name} {distributorUser.last_name}
                                </p>
                            </div>

                            <div className="space-y-4">
                                <div className="grid grid-cols-1 gap-6">
                                    <div className="space-y-4">
                                        <div className="grid gap-2">
                                            <Label htmlFor="FirstName">First Name</Label>
                                            <Input
                                                id="FirstName"
                                                name="FirstName"
                                                defaultValue={distributorUser.first_name}
                                            />
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="LastName">Last Name</Label>
                                            <Input
                                                id="LastName"
                                                name="LastName"
                                                defaultValue={distributorUser.last_name}
                                                required
                                            />
                                        </div>

                                        <div className="grid gap-2">
                                            <Label htmlFor="Email">Email</Label>
                                            <Input
                                                id="Email"
                                                name="Email"
                                                defaultValue={distributorUser.email}
                                                type="email"
                                                placeholder="abebe.kebede@example.com"
                                                required
                                            />
                                        </div>

                                        <div className="grid gap-2">
                                            <Label htmlFor="phone">Phone</Label>
                                            <Input
                                                id="phone"
                                                name="phone"
                                                defaultValue={distributorUser.phone}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value="business">
                    <Card className="rounded-sm border-2 border-gray-200 shadow-none">
                        <CardContent className="">
                            <div className="mb-4">
                                <h3 className="text-xl font-bold">
                                    Business Information
                                </h3>
                                <p className="text-sm text-muted-foreground">
                                    Add Business Information
                                </p>
                            </div>

                            <div className="space-y-4">
                                <div className="grid grid-cols-1 gap-6">
                                    <div className="space-y-4">
                                        <div className="grid gap-2">
                                            <Label htmlFor="Name">Name</Label>
                                            <Input
                                                id="Name"
                                                name="Name"
                                                defaultValue={distributor.name}
                                                placeholder="Business Name"
                                            />
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="Tin">Tin</Label>
                                            <Input
                                                id="Tin"
                                                name="Tin"
                                                defaultValue={distributor.tin}

                                            />
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="General Zone">General Zone</Label>
                                            <Input
                                                id="General Zone"
                                                name="General Zone"
                                                defaultValue={distributor.general_zone}
                                            />
                                        </div>

                                        <div className="grid gap-2">
                                            <Label htmlFor="Region">Region</Label>
                                            <Input
                                                id="Region"
                                                name="Region"
                                                defaultValue={distributor.region}

                                            />
                                        </div>

                                        <div className="grid gap-2">
                                            <Label htmlFor="Woreda">Woreda</Label>
                                            <Input
                                                id="Woreda"
                                                name="Woreda"
                                                defaultValue={distributor.woreda}

                                            />
                                        </div>

                                        {
                                            distributor && (
                                                <div className="h-[400px] rounded-lg overflow-hidden">
                                                    <Map
                                                        markerPosition={markerPosition}
                                                        useCurrentLocation={true}
                                                        setMarkerPosition={setMarkerPosition}
                                                        setFormData={setFormData}
                                                    />
                                                </div>
                                            )
                                        }

                                    </div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
            <div className="flex items-center justify-end space-x-4">
                <Button type="submit" onClick={handleSubmit}>
                    Approve Distributor
                </Button>
            </div>
        </div>
    );
}
