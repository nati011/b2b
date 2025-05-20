'use client'
import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight, Plus, X } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";

export interface Category {
    id: number;
    name: string;
}

interface CategoriesFormProps {
    categories: Category[];
    selectedCategories: number[];
    onSelectCategories: (categoryIds: number[]) => void;
    onAddCategory: (categoryName: string) => void;
    onDeleteCategory: (categoryId: number) => void;
    onBack: () => void;
    onNext: () => void;
}

const CategoriesForm: React.FC<CategoriesFormProps> = ({
    categories,
    selectedCategories,
    onSelectCategories,
    onAddCategory,
    onDeleteCategory,
    onBack,
    onNext,
}) => {
    const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [newCategoryName, setNewCategoryName] = useState("");
    const [categoryToDelete, setCategoryToDelete] = useState<number | null>(null);
    const [searchTerm, setSearchTerm] = useState("");


    const filteredCategories = categories.filter(category =>
        category.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const toggleCategory = (categoryId: number) => {
        if (selectedCategories.includes(categoryId)) {
            onSelectCategories(selectedCategories.filter(id => id !== categoryId));
        } else {
            onSelectCategories([...selectedCategories, categoryId]);
        }
    };

    const handleAddCategory = () => {
        if (newCategoryName.trim()) {
            onAddCategory(newCategoryName.trim());
            setNewCategoryName("");
            setIsAddDialogOpen(false);
        }
    };

    const confirmDeleteCategory = (categoryId: number) => {
        setCategoryToDelete(categoryId);
        setIsDeleteDialogOpen(true);
    };

    const handleDeleteCategory = () => {
        if (categoryToDelete !== null) {
            onDeleteCategory(categoryToDelete);
            setIsDeleteDialogOpen(false);
            setCategoryToDelete(null);
        }
    };

    return (
        <div className="space-y-6">
            <div className="space-y-4">
                <div className="flex justify-between items-center">
                    <h3 className="text-lg font-medium">Product Categories</h3>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setIsAddDialogOpen(true)}
                        className="flex items-center gap-1"
                    >
                        <Plus className="w-4 h-4" /> Add Category
                    </Button>
                </div>

                <div className="space-y-4">
                    <Input
                        placeholder="Search categories..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="mb-2"
                    />

                    <ScrollArea className="h-[240px] border rounded-md p-4">
                        {filteredCategories.length === 0 ? (
                            <p className="text-center text-muted-foreground p-4">No categories found</p>
                        ) : (
                            <div className="space-y-2">
                                {filteredCategories.map((category) => (
                                    <div key={category.id} className="flex items-center justify-between p-2 hover:bg-accent rounded-md">
                                        <div
                                            className="flex-1 cursor-pointer flex items-center"
                                            onClick={() => toggleCategory(category.id)}
                                        >
                                            <div className={`w-4 h-4 border rounded mr-2 flex items-center justify-center ${selectedCategories.includes(category.id) ? "bg-primary border-primary" : ""
                                                }`}>
                                                {selectedCategories.includes(category.id) && (
                                                    <svg width="10" height="8" viewBox="0 0 10 8" fill="none" xmlns="http://www.w3.org/2000/svg">
                                                        <path d="M9 1L3.5 6.5L1 4" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                                                    </svg>
                                                )}
                                            </div>
                                            <span>{category.name}</span>
                                        </div>
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                confirmDeleteCategory(category.id);
                                            }}
                                        >
                                            <X className="w-4 h-4 text-muted-foreground hover:text-destructive" />
                                        </Button>
                                    </div>
                                ))}
                            </div>
                        )}
                    </ScrollArea>

                    <div className="flex flex-wrap gap-2 mt-4">
                        {selectedCategories.map(categoryId => {
                            const category = categories.find(c => c.id === categoryId);
                            return category ? (
                                <Badge key={categoryId} variant="outline" className="px-3 py-1">
                                    {category.name}
                                    <X
                                        className="ml-2 w-3 h-3 cursor-pointer"
                                        onClick={() => toggleCategory(categoryId)}
                                    />
                                </Badge>
                            ) : null;
                        })}
                        {selectedCategories.length === 0 && (
                            <p className="text-sm text-muted-foreground italic">No categories selected</p>
                        )}
                    </div>
                </div>
            </div>

            <div className="flex justify-between">
                <Button type="button" variant="outline" onClick={onBack} className="flex items-center gap-2">
                    <ArrowLeft className="w-4 h-4" /> Back
                </Button>
                <Button onClick={onNext} className="flex items-center gap-2">
                    Continue <ArrowRight className="w-4 h-4" />
                </Button>
            </div>

            {/* Add Category Dialog */}
            <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Add New Category</DialogTitle>
                    </DialogHeader>
                    <Input
                        placeholder="Category Name"
                        value={newCategoryName}
                        onChange={(e) => setNewCategoryName(e.target.value)}
                        autoFocus
                    />
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsAddDialogOpen(false)}>Cancel</Button>
                        <Button onClick={handleAddCategory}>Add Category</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Delete Confirmation Dialog */}
            <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Confirm Delete</DialogTitle>
                    </DialogHeader>
                    <p>Are you sure you want to delete this category?</p>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>Cancel</Button>
                        <Button variant="destructive" onClick={handleDeleteCategory}>Delete</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
};

export default CategoriesForm;