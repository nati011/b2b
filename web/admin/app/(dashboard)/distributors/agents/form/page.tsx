'use client'
import { useEffect, useState } from "react";
import useDistributorsStore from "@/app/libs/store/useDistributorStore";
import { UserPlus, ArrowLeft, Save, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import Breadcrumb from "@/components/breadcrumb";
import Link from "next/link";
import { useRouter } from "next/navigation";

interface FormData {
    first_name: string;
    last_name: string;
    username: string;
    email: string;
    phone: string;
}

interface FormErrors {
    first_name?: string;
    last_name?: string;
    username?: string;
    email?: string;
    phone?: string;
}

export default function AgentForm() {
    const router = useRouter();
    const [formData, setFormData] = useState<FormData>({
        first_name: '',
        last_name: '',
        username: '',
        email: '',
        phone: ''
    });
    const [errors, setErrors] = useState<FormErrors>({});
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [submitError, setSubmitError] = useState<string>('');
    const {
        createDistirbutorUser
    } = useDistributorsStore()

    const validateForm = (): boolean => {
        const newErrors: FormErrors = {};

        // First Name validation
        if (!formData.first_name.trim()) {
            newErrors.first_name = 'First name is required';
        } else if (formData.first_name.trim().length < 2) {
            newErrors.first_name = 'First name must be at least 2 characters';
        }

        // Last Name validation
        if (!formData.last_name.trim()) {
            newErrors.last_name = 'Last name is required';
        } else if (formData.last_name.trim().length < 2) {
            newErrors.last_name = 'Last name must be at least 2 characters';
        }

        // Username validation
        if (!formData.username.trim()) {
            newErrors.username = 'Username is required';
        } else if (formData.username.trim().length < 3) {
            newErrors.username = 'Username must be at least 3 characters';
        } else if (!/^[a-zA-Z0-9_]+$/.test(formData.username)) {
            newErrors.username = 'Username can only contain letters, numbers, and underscores';
        }

        // Email validation
        if (!formData.email.trim()) {
            newErrors.email = 'Email is required';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
            newErrors.email = 'Please enter a valid email address';
        }

        // Phone validation
        if (!formData.phone.trim()) {
            newErrors.phone = 'Phone number is required';
        } else if (!/^[\+]?[1-9][\d]{0,15}$/.test(formData.phone.replace(/\s/g, ''))) {
            newErrors.phone = 'Please enter a valid phone number';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleInputChange = (field: keyof FormData, value: string) => {
        setFormData(prev => ({
            ...prev,
            [field]: value
        }));
   
        if (errors[field]) {
            setErrors(prev => ({
                ...prev,
                [field]: undefined
            }));
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitError('');
        if (!validateForm()) {
            return;
        }
        setIsSubmitting(true);
        try {
            createDistirbutorUser(formData)
            router.push('/distributors/agents');
        } catch (error) {
            setSubmitError('Failed to create agent. Please try again.');
            console.error('Error creating agent:', error);
        } finally {
            setIsSubmitting(false);
        }
    };

    const breadcrumbItems = [
        { title: 'Distributors', href: '/distributors' },
        { title: 'Agents', href: '/distributors/agents' },
        { title: 'Create Agent', href: '/distributors/agents/form' }
    ];

    return (
        <div className="space-y-6">
            <Breadcrumb 
                page={breadcrumbItems}
                heading="Create New Agent"
                subheading="Add a new distributor agent to the system"
            />

            <Card className="">
                <CardHeader className="border-b ">
                    <CardTitle className="flex items-center gap-2">
                        <UserPlus className="h-5 w-5" />
                        Agent Information
                    </CardTitle>
                    <CardDescription>
                        Fill in the details below to create a new distributor agent account.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {submitError && (
                            <Alert variant="destructive">
                                <AlertDescription>{submitError}</AlertDescription>
                            </Alert>
                        )}

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="first_name">First Name *</Label>
                                <Input
                                    id="first_name"
                                    type="text"
                                    placeholder="Enter first name"
                                    value={formData.first_name}
                                    onChange={(e) => handleInputChange('first_name', e.target.value)}
                                    className={errors.first_name ? 'border-destructive' : ''}
                                />
                                {errors.first_name && (
                                    <p className="text-sm text-destructive">{errors.first_name}</p>
                                )}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="last_name">Last Name *</Label>
                                <Input
                                    id="last_name"
                                    type="text"
                                    placeholder="Enter last name"
                                    value={formData.last_name}
                                    onChange={(e) => handleInputChange('last_name', e.target.value)}
                                    className={errors.last_name ? 'border-destructive' : ''}
                                />
                                {errors.last_name && (
                                    <p className="text-sm text-destructive">{errors.last_name}</p>
                                )}
                            </div>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="username">Username *</Label>
                            <Input
                                id="username"
                                type="text"
                                placeholder="Enter username"
                                value={formData.username}
                                onChange={(e) => handleInputChange('username', e.target.value)}
                                className={errors.username ? 'border-destructive' : ''}
                            />
                            {errors.username && (
                                <p className="text-sm text-destructive">{errors.username}</p>
                            )}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="email">Email Address *</Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="Enter email address"
                                value={formData.email}
                                onChange={(e) => handleInputChange('email', e.target.value)}
                                className={errors.email ? 'border-destructive' : ''}
                            />
                            {errors.email && (
                                <p className="text-sm text-destructive">{errors.email}</p>
                            )}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="phone">Phone Number *</Label>
                            <Input
                                id="phone"
                                type="tel"
                                placeholder="Enter phone number"
                                value={formData.phone}
                                onChange={(e) => handleInputChange('phone', e.target.value)}
                                className={errors.phone ? 'border-destructive' : ''}
                            />
                            {errors.phone && (
                                <p className="text-sm text-destructive">{errors.phone}</p>
                            )}
                        </div>

                        <div className="flex items-center justify-end gap-4 pt-4">
                            <Link href="/distributors/agents">
                                <Button type="button" variant="outline">
                                    Cancel
                                </Button>
                            </Link>
                            <Button type="submit" disabled={isSubmitting} className="gap-2">
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="h-4 w-4 animate-spin" />
                                        Creating...
                                    </>
                                ) : (
                                    <>
                                        <Save className="h-4 w-4" />
                                        Create Agent
                                    </>
                                )}
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
}