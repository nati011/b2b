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
import { Checkbox } from "@/components/ui/checkbox";
import { useRouter } from "next/navigation";

const Map = dynamic(
  () => import('@/app/components/map'),
  { ssr: false }
)


export default function DistributorsForm() {
  const {
    loading,
    error,
    createDistributors
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
  const [useCurrentLocation, setUseCurrentLocation] = useState(false)
  const [submitted, setSubmitted] = useState(false)
  const router = useRouter()

  const handleReset = () => {
    setFormData({
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

    });
  };


  useEffect(() => {
    if (useCurrentLocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const newPosition: [number, number] = [
            position.coords.latitude,
            position.coords.longitude
          ]
          console.log(position.coords.latitude,
            position.coords.longitude)
          setMarkerPosition(newPosition)
          setFormData(prev => ({ ...prev, latitude: newPosition[0].toString(), longitude: newPosition[1].toString() }))
        },
        (error) => {
          console.error('Error getting location:', error)
        }
      )
    }
  }, [useCurrentLocation])

  useEffect(() => {
    if (submitted && !loading) {
      if (error) {
        toast.error(error)
      } else {
        toast.success("Distributor created successfully!")
        handleReset()
        router.push('/distributors')
      }
      setSubmitted(false)
    }
  }, [loading, error, submitted])
  const pages = [
    {
      "title": "Distributor",
      "href": "/distributor"
    },
    {
      "title": "Form",
      "href": "/distributor/form"
    }
  ]

  const handleSubmit = () => {
    setSubmitted(true)
    createDistributors(formData)
  }

  return (
    <div className="grid grid-cols-1 gap-4">
      <Heading page={pages} heading="Register Distributors" subheading="Fill in the forms accorfingly to register." />
      <Card className="rounded-sm border-2 border-gray-200 shadow-none">
        <CardContent className="">
          <div className="mb-4">
            <h3 className="text-xl font-bold">
              Profile Information
            </h3>
            <p className="text-sm text-muted-foreground">
              Register a new distributor
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
                    value={formData.first_name}
                    onChange={(e) => { setFormData(prev => ({ ...prev, first_name: e.target.value })) }}
                    placeholder="Abebe"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="LastName">Last Name</Label>
                  <Input
                    id="LastName"
                    name="LastName"
                    value={formData.last_name}
                    onChange={(e) => { setFormData(prev => ({ ...prev, last_name: e.target.value })) }}
                    placeholder="Kebede"
                    required
                  />
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="Email">Email</Label>
                  <Input
                    id="Email"
                    name="Email"
                    value={formData.email}
                    type="email"
                    onChange={(e) => { setFormData(prev => ({ ...prev, email: e.target.value })) }}
                    placeholder="abebe.kebede@example.com"
                    required
                  />
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="phone">Phone</Label>
                  <Input
                    id="phone"
                    name="phone"
                    value={formData.phone}
                    onChange={(e) => { setFormData(prev => ({ ...prev, phone: e.target.value })) }}
                    placeholder="+251955123456"
                    required
                  />
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
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
                    value={formData.name}
                    onChange={(e) => { setFormData(prev => ({ ...prev, name: e.target.value })) }}
                    placeholder="Business Name"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="Tin">Tin</Label>
                  <Input
                    id="Tin"
                    name="Tin"
                    value={formData.tin}
                    type="number"
                    onChange={(e) => { setFormData(prev => ({ ...prev, tin: e.target.value })) }}
                    placeholder="1234567890"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="General Zone">General Zone</Label>
                  <Input
                    id="General Zone"
                    name="General Zone"
                    value={formData.general_zone}
                    onChange={(e) => { setFormData(prev => ({ ...prev, general_zone: e.target.value })) }}
                    placeholder="Bole"
                    required
                  />
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="Region">Region</Label>
                  <Input
                    id="Region"
                    name="Region"
                    value={formData.region}
                    onChange={(e) => { setFormData(prev => ({ ...prev, region: e.target.value })) }}
                    placeholder="region-001"
                    required
                  />
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="Woreda">Woreda</Label>
                  <Input
                    id="Woreda"
                    name="Woreda"
                    value={formData.woreda}
                    onChange={(e) => { setFormData(prev => ({ ...prev, woreda: e.target.value })) }}
                    placeholder="woreda-001"
                    required
                  />
                </div>
                <div className="flex items-center space-x-2 my-4">
                  <Checkbox
                    id="useLocation"
                    checked={useCurrentLocation}
                    onCheckedChange={(checked: any) => setUseCurrentLocation(checked as boolean)}
                  />
                  <label htmlFor="useLocation">Use my current location</label>
                </div>

                <div className="h-[400px] rounded-lg overflow-hidden">
                  <Map
                    markerPosition={markerPosition}
                    useCurrentLocation={useCurrentLocation}
                    setMarkerPosition={setMarkerPosition}
                    setFormData={setFormData}
                  />
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>




      <div className="flex items-center justify-end space-x-4">
        <Button
          type="button"
          variant="outline"
          onClick={handleReset}
          disabled={loading}
        >
          Reset
        </Button>
        <Button
          type="submit"
          onClick={handleSubmit}
          disabled={loading}
        >
          {loading ? "Registering..." : "Register Distributor"}
        </Button>
      </div>
    </div>
  );
}
