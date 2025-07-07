export const Footer = () => {
    const currentYear = new Date().getFullYear();

    return (
        <div className="mt-16 w-full print:hidden">
            <div className="container mx-auto px-4 py-12">
                <div className="border-t border-primary mt-12 pt-6">
                    <div className="flex flex-col md:flex-row justify-between items-center">
                        <p className="text-sm text-muted-foreground mb-4 md:mb-0">
                            &copy; {currentYear} Efoyeta. All rights reserved.
                        </p>
                        <div className="flex space-x-6">
                            <a href="#" className="text-xs text-muted-foreground hover:text-primary transition-colors">
                                Privacy Policy
                            </a>
                            <a href="#" className="text-xs text-muted-foreground hover:text-primary transition-colors">
                                Terms of Service
                            </a>
                            <a href="#" className="text-xs text-muted-foreground hover:text-primary transition-colors">
                                Cookies
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};