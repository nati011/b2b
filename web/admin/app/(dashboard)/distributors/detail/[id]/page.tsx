'use client'
import Heading from "@/components/breadcrumb";
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { Distributor, DistributorRequest, DistributorVerdict } from '@/app/libs/types';
import dynamic from "next/dynamic";
import 'leaflet/dist/leaflet.css'
import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { useParams, useRouter } from "next/navigation";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { CheckCircle, XCircle, Clock, PlayCircle, PauseCircle, MapPin, Building2, User } from "lucide-react";
import { useUserStore } from "@/app/libs/store/useAuthStore";

const Map = dynamic(
    () => import('@/components/map'),
    { ssr: false }
)

export default function DistributorDetail() {
    const {
        success,
        distributor,
        distributorUser,
        loading,
        error,
        fetchDistributorDetail,
        approveDistributor,
        rejectDistributor,
        activateDistributor,
        deactivateDistributor
    } = useDistributorsStore()

    const {
        user,
        fetchUserById
    } = useUserStore()
    
    const [markerPosition, setMarkerPosition] = useState<[number, number]>([8.9934609, 38.7714897])
    const [comment, setComment] = useState("")
    const [approvalDialogOpen, setApprovalDialogOpen] = useState(false)
    const [rejectionDialogOpen, setRejectionDialogOpen] = useState(false)
    const [deactivateDialogOpen, setDeactivateDialogOpen] = useState(false)
    const [actionType, setActionType] = useState<'approve' | 'reject' | 'activate' | 'deactivate'>('approve')
    
    const router = useRouter()
    const routeParam = useParams<{ id: string }>();
    
    const pages = [
        {
            "title": "Distributors",
            "href": "/distributors"
        },
        {
            "title": "Detail",
            "href": `/distributors/detail/${routeParam.id}`
        }
    ]

    useEffect(() => {
        if (routeParam.id) {
            fetchDistributorDetail(parseInt(routeParam.id))
        }
    }, [routeParam.id])

    useEffect(() => {
        if (distributor?.user && distributor?.user[0] !== 0) {
            fetchUserById(distributor?.user[0])
            const lat = distributor?.latitude || "8.9934609"
            const long = distributor?.longitude || "38.7714897"
            setMarkerPosition([parseFloat(lat), parseFloat(long)])
        }
    }, [distributor])

    useEffect(() => {
        if (error) {
            toast.error(error)
        }
        if (success) {
            toast.success(success)
        }
    }, [error, success])

    const getStatusBadge = () => {
        if (distributor?.verdict === DistributorVerdict.PENDING) {
            return (
                <Badge variant="secondary" className="flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    Pending Approval
                </Badge>
            );
        }
        if (distributor?.verdict === DistributorVerdict.REJECTED) {
            return (
                <Badge variant="destructive" className="flex items-center gap-1">
                    <XCircle className="w-3 h-3" />
                    Rejected
                </Badge>
            );
        }
        if (distributor?.verdict === DistributorVerdict.APPROVED && distributor?.is_active) {
            return (
                <Badge variant="default" className="flex items-center gap-1 bg-green-100 text-green-800 hover:bg-green-100">
                    <CheckCircle className="w-3 h-3" />
                    Active
                </Badge>
            );
        }
        if (distributor?.verdict === DistributorVerdict.APPROVED && !distributor?.is_active) {
            return (
                <Badge variant="destructive" className="flex items-center gap-1">
                    <PauseCircle className="w-3 h-3" />
                    Inactive
                </Badge>
            );
        }
        return null;
    };

    const handleAction = (type: 'approve' | 'reject' | 'activate' | 'deactivate') => {
        setActionType(type)
        setComment("")
        
        if (type === 'approve') {
            setApprovalDialogOpen(true)
        } else if (type === 'reject') {
            setRejectionDialogOpen(true)
        } else if (type === 'deactivate') {
            setDeactivateDialogOpen(true)
        } else {
            // Direct activation
            activateDistributor(parseInt(routeParam.id))
        }
    }

    const handleConfirmAction = async () => {
        try {
            if (actionType === 'approve') {
                await approveDistributor(parseInt(routeParam.id))
                setApprovalDialogOpen(false)
            } else if (actionType === 'reject') {
                if (!comment.trim()) {
                    toast.error("Please provide a reason for rejection")
                    return
                }
                await rejectDistributor(parseInt(routeParam.id), comment)
                setRejectionDialogOpen(false)
            } else if (actionType === 'deactivate') {
                await deactivateDistributor(parseInt(routeParam.id))
                setDeactivateDialogOpen(false)
            }
        } catch (error) {
            console.error('Action failed:', error)
        }
    }

    if (loading && !distributor?.id) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
                    <p className="mt-2 text-sm text-gray-600">Loading distributor details...</p>
                </div>
            </div>
        )
    }

    return (
        <div className="space-y-6">
                <Heading page={pages} heading="Distributor Details" subheading={`${distributor?.name || 'Loading...'}`} />
                <div className="ml-auto w-full bg-red">
                    {getStatusBadge()}
                </div>

            <Tabs defaultValue="business" className="w-full">
                <TabsList className="w-full">
                    <TabsTrigger value="business" className="flex items-center gap-2">
                        <Building2 className="w-4 h-4" />
                        Business Information
                    </TabsTrigger>
                    <TabsTrigger value="profile" className="flex items-center gap-2">
                        <User className="w-4 h-4" />
                        Profile Information
                    </TabsTrigger>
                    <TabsTrigger value="location" className="flex items-center gap-2">
                        <MapPin className="w-4 h-4" />
                        Location
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="business">
                    <Card className="rounded-sm">
                        <CardHeader>
                            <CardTitle>Business Information</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div className="space-y-2">
                                    <Label htmlFor="name">Business Name</Label>
                                    <Input
                                        id="name"
                                        value={distributor?.name || ''}
                                        readOnly
                                        className="bg-gray-50"
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label htmlFor="tin">TIN Number</Label>
                                    <Input
                                        id="tin"
                                        value={distributor?.tin || ''}
                                        readOnly
                                        className="bg-gray-50 font-mono"
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label htmlFor="region">Region</Label>
                                    <Input
                                        id="region"
                                        value={distributor?.region || ''}
                                        readOnly
                                        className="bg-gray-50"
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label htmlFor="zone">General Zone</Label>
                                    <Input
                                        id="zone"
                                        value={distributor?.general_zone || ''}
                                        readOnly
                                        className="bg-gray-50"
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label htmlFor="woreda">Woreda</Label>
                                    <Input
                                        id="woreda"
                                        value={distributor?.woreda || ''}
                                        readOnly
                                        className="bg-gray-50"
                                    />
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="profile">
                    <Card className="rounded-sm">
                        <CardHeader>
                            <CardTitle>Profile Information</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            {user?.id ? (
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <Label htmlFor="firstName">First Name</Label>
                                        <Input
                                            id="firstName"
                                            value={user?.first_name || ''}
                                            readOnly
                                            className="bg-gray-50"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="lastName">Last Name</Label>
                                        <Input
                                            id="lastName"
                                            value={user?.last_name || ''}
                                            readOnly
                                            className="bg-gray-50"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="email">Email</Label>
                                        <Input
                                            id="email"
                                            value={user?.email || ''}
                                            readOnly
                                            className="bg-gray-50"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="phone">Phone</Label>
                                        <Input
                                            id="phone"
                                            value={user?.phone || ''}
                                            readOnly
                                            className="bg-gray-50"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="username">Username</Label>
                                        <Input
                                            id="username"
                                            value={user?.username || ''}
                                            readOnly
                                            className="bg-gray-50"
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="dob">Date of Birth</Label>
                                        <Input
                                            id="dob"
                                            value={user?.dob || ''}
                                            readOnly
                                            className="bg-gray-50"
                                        />
                                    </div>
                                </div>
                            ) : (
                                <div className="text-center py-8 text-gray-500">
                                    <User className="w-12 h-12 mx-auto mb-4 opacity-50" />
                                    <p>No user profile information available</p>
                                </div>
                            )}
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="location">
                    <Card className="rounded-sm">
                        <CardHeader>
                            <CardTitle>Location Information</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                                <div className="space-y-2">
                                    <Label htmlFor="latitude">Latitude</Label>
                                    <Input
                                        id="latitude"
                                        value={distributor?.latitude || ''}
                                        readOnly
                                        className="bg-gray-50 font-mono"
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label htmlFor="longitude">Longitude</Label>
                                    <Input
                                        id="longitude"
                                        value={distributor?.longitude || ''}
                                        readOnly
                                        className="bg-gray-50 font-mono"
                                    />
                                </div>
                            </div>
                            <div className="h-[400px] rounded-lg overflow-hidden border">
                                <Map
                                    markerPosition={markerPosition}
                                    useCurrentLocation={false}
                                    setMarkerPosition={setMarkerPosition}
                                    setFormData={() => {}}
                                />
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>

            {/* Action Buttons */}
            <div className="flex items-center justify-end gap-3">
                {distributor?.verdict === DistributorVerdict.PENDING ? (
                    <>
                        <Button 
                            onClick={() => handleAction('approve')} 
                            disabled={loading}
                            variant="outline"
                            className="border-green-600 text-green-600 hover:bg-green-50"
                        >
                            <CheckCircle className="w-4 h-4 mr-2" />
                            {loading ? "Processing..." : "Approve"}
                        </Button>
                        <Button 
                            onClick={() => handleAction('reject')} 
                            disabled={loading}
                            variant="outline"
                            className="border-red-600 text-red-600 hover:bg-red-50"
                        >
                            <XCircle className="w-4 h-4 mr-2" />
                            {loading ? "Processing..." : "Reject"}
                        </Button>
                    </>
                ) : (
                    <>
                        {distributor?.verdict === DistributorVerdict.APPROVED && !distributor?.is_active ? (
                            <Button 
                                onClick={() => handleAction('activate')} 
                                disabled={loading}
                                variant="outline"
                                className="border-blue-600 text-blue-600 hover:bg-blue-50"
                            >
                                <PlayCircle className="w-4 h-4 mr-2" />
                                {loading ? "Processing..." : "Activate"}
                            </Button>
                        ) : null}
                        {distributor?.verdict === DistributorVerdict.APPROVED && distributor?.is_active ? (
                            <Button 
                                onClick={() => handleAction('deactivate')} 
                                disabled={loading}
                                variant="outline"
                                className="border-orange-600 text-orange-600 hover:bg-orange-50"
                            >
                                <PauseCircle className="w-4 h-4 mr-2" />
                                {loading ? "Processing..." : "Deactivate"}
                            </Button>
                        ) : null}
                    </>
                )}
            </div>

            {/* Approval Dialog */}
            <Dialog open={approvalDialogOpen} onOpenChange={setApprovalDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-green-600" />
                            Confirm Approval
                        </DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <div className="space-y-4">
                        <p className="text-sm text-gray-600">
                            Are you sure you want to approve this distributor? This action will allow them to access the platform.
                        </p>
                        <div className="grid gap-2">
                            <Label htmlFor="approvalComment">Optional Comment</Label>
                            <Textarea
                                id="approvalComment"
                                value={comment}
                                onChange={(e) => setComment(e.target.value)}
                                placeholder="Add any notes about this approval..."
                                rows={3}
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setApprovalDialogOpen(false)}>
                            Cancel
                        </Button>
                        <Button onClick={handleConfirmAction} className="bg-green-600 hover:bg-green-700">
                            Approve Distributor
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Rejection Dialog */}
            <Dialog open={rejectionDialogOpen} onOpenChange={setRejectionDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <XCircle className="w-5 h-5 text-red-600" />
                            Confirm Rejection
                        </DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <div className="space-y-4">
                        <p className="text-sm text-gray-600">
                            Please provide a reason for rejecting this distributor application.
                        </p>
                        <div className="grid gap-2">
                            <Label htmlFor="rejectionComment">Reason for Rejection *</Label>
                            <Textarea
                                id="rejectionComment"
                                value={comment}
                                onChange={(e) => setComment(e.target.value)}
                                placeholder="Explain why this application is being rejected..."
                                rows={4}
                                required
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setRejectionDialogOpen(false)}>
                            Cancel
                        </Button>
                        <Button 
                            onClick={handleConfirmAction} 
                            variant="destructive"
                            disabled={!comment.trim()}
                        >
                            Reject Application
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Deactivation Dialog */}
            <Dialog open={deactivateDialogOpen} onOpenChange={setDeactivateDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <PauseCircle className="w-5 h-5 text-orange-600" />
                            Confirm Deactivation
                        </DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <div className="space-y-4">
                        <p className="text-sm text-gray-600">
                            Are you sure you want to deactivate this distributor? They will lose access to the platform until reactivated.
                        </p>
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setDeactivateDialogOpen(false)}>
                            Cancel
                        </Button>
                        <Button onClick={handleConfirmAction} variant="destructive">
                            Deactivate Distributor
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}
