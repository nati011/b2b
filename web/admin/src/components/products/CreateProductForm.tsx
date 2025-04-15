
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";

type ProductFormData = {
  name: string;
  category: string;
  price: string;
  stock: string;
  description: string;
  features: string[];
};

const initialFormData: ProductFormData = {
  name: "",
  category: "",
  price: "",
  stock: "",
  description: "",
  features: [""],
};

const categories = ["Electronics", "Tools", "Appliances", "Clothing", "Furniture"];

const CreateProductForm = ({ 
  open, 
  onOpenChange 
}: { 
  open: boolean; 
  onOpenChange: (open: boolean) => void;
}) => {
  const [step, setStep] = useState(1);
  const [formData, setFormData] = useState<ProductFormData>(initialFormData);
  const [errors, setErrors] = useState<Partial<Record<keyof ProductFormData, string>>>({});
  const { toast } = useToast();
  
  const totalSteps = 3;

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
    // Clear error when field is edited
    if (errors[name as keyof ProductFormData]) {
      setErrors({ ...errors, [name]: "" });
    }
  };

  const handleSelectChange = (value: string, name: string) => {
    setFormData({ ...formData, [name]: value });
    // Clear error when field is edited
    if (errors[name as keyof ProductFormData]) {
      setErrors({ ...errors, [name]: "" });
    }
  };

  const handleFeatureChange = (index: number, value: string) => {
    const newFeatures = [...formData.features];
    newFeatures[index] = value;
    setFormData({ ...formData, features: newFeatures });
    
    // Clear feature errors if they exist
    if (errors.features) {
      setErrors({ ...errors, features: "" });
    }
  };

  const addFeature = () => {
    setFormData({ ...formData, features: [...formData.features, ""] });
  };

  const removeFeature = (index: number) => {
    const newFeatures = formData.features.filter((_, i) => i !== index);
    setFormData({ ...formData, features: newFeatures });
  };

  const validateStep = (currentStep: number): boolean => {
    const newErrors: Partial<Record<keyof ProductFormData, string>> = {};
    let isValid = true;

    if (currentStep === 1) {
      if (!formData.name.trim()) {
        newErrors.name = "Product name is required";
        isValid = false;
      }
      if (!formData.category) {
        newErrors.category = "Category is required";
        isValid = false;
      }
    } else if (currentStep === 2) {
      if (!formData.price.trim()) {
        newErrors.price = "Price is required";
        isValid = false;
      } else if (isNaN(parseFloat(formData.price)) || parseFloat(formData.price) <= 0) {
        newErrors.price = "Price must be a positive number";
        isValid = false;
      }
      
      if (!formData.stock.trim()) {
        newErrors.stock = "Stock quantity is required";
        isValid = false;
      } else if (isNaN(parseInt(formData.stock)) || parseInt(formData.stock) < 0) {
        newErrors.stock = "Stock must be a non-negative number";
        isValid = false;
      }
    } else if (currentStep === 3) {
      const emptyFeatures = formData.features.filter(f => !f.trim()).length;
      if (emptyFeatures > 0) {
        newErrors.features = "All features must have a value or be removed";
        isValid = false;
      }
    }

    setErrors(newErrors);
    return isValid;
  };

  const nextStep = () => {
    if (validateStep(step) && step < totalSteps) {
      setStep(step + 1);
    }
  };

  const prevStep = () => {
    if (step > 1) {
      setStep(step - 1);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validateStep(step)) {
      // Here you would typically handle form submission
      console.log("Product form submitted:", formData);
      toast({
        title: "Product Created",
        description: `${formData.name} has been added to your inventory.`,
      });
      setFormData(initialFormData);
      setStep(1);
      onOpenChange(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Create New Product</DialogTitle>
          <DialogDescription>
            Step {step} of {totalSteps}: {step === 1 ? "Basic Info" : step === 2 ? "Details" : "Features"}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          {step === 1 && (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Product Name</Label>
                <Input
                  id="name"
                  name="name"
                  placeholder="Enter product name"
                  value={formData.name}
                  onChange={handleInputChange}
                  className={errors.name ? "border-red-500" : ""}
                />
                {errors.name && <p className="text-sm text-red-500">{errors.name}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="category">Category</Label>
                <Select
                  value={formData.category}
                  onValueChange={(value) => handleSelectChange(value, "category")}
                >
                  <SelectTrigger className={errors.category ? "border-red-500" : ""}>
                    <SelectValue placeholder="Select category" />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map((category) => (
                      <SelectItem key={category} value={category}>
                        {category}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.category && <p className="text-sm text-red-500">{errors.category}</p>}
              </div>
            </div>
          )}

          {step === 2 && (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="price">Price ($)</Label>
                <Input
                  id="price"
                  name="price"
                  placeholder="0.00"
                  value={formData.price}
                  onChange={handleInputChange}
                  className={errors.price ? "border-red-500" : ""}
                />
                {errors.price && <p className="text-sm text-red-500">{errors.price}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="stock">Stock Quantity</Label>
                <Input
                  id="stock"
                  name="stock"
                  type="number"
                  placeholder="0"
                  value={formData.stock}
                  onChange={handleInputChange}
                  className={errors.stock ? "border-red-500" : ""}
                />
                {errors.stock && <p className="text-sm text-red-500">{errors.stock}</p>}
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <textarea
                  id="description"
                  name="description"
                  rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-base ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
                  placeholder="Enter product description"
                  value={formData.description}
                  onChange={handleInputChange}
                />
              </div>
            </div>
          )}

          {step === 3 && (
            <div className="space-y-4">
              <div>
                <Label>Product Features</Label>
                <p className="text-sm text-muted-foreground mb-2">
                  Add key features of your product
                </p>
                {formData.features.map((feature, index) => (
                  <div key={index} className="flex items-center gap-2 mb-2">
                    <Input
                      placeholder={`Feature ${index + 1}`}
                      value={feature}
                      onChange={(e) => handleFeatureChange(index, e.target.value)}
                      className={errors.features ? "border-red-500" : ""}
                    />
                    {index > 0 && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => removeFeature(index)}
                      >
                        Remove
                      </Button>
                    )}
                  </div>
                ))}
                {errors.features && <p className="text-sm text-red-500">{errors.features}</p>}
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={addFeature}
                  className="mt-2"
                >
                  Add Feature
                </Button>
              </div>
            </div>
          )}

          <DialogFooter className="mt-6">
            {step > 1 && (
              <Button type="button" variant="outline" onClick={prevStep}>
                Back
              </Button>
            )}
            {step < totalSteps ? (
              <Button type="button" onClick={nextStep}>
                Next
              </Button>
            ) : (
              <Button type="submit">Create Product</Button>
            )}
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default CreateProductForm;
