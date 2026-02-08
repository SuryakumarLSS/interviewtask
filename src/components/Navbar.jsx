import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Menu, X } from 'lucide-react';
import { cn } from '../lib/utils';

const Navbar = () => {
    const [isScrolled, setIsScrolled] = useState(false);
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            setIsScrolled(window.scrollY > 20);
        };
        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, []);

    return (
        <>
            <motion.nav
                initial={{ y: -100, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ duration: 0.5 }}
                className={cn(
                    "fixed top-0 left-0 right-0 z-50 flex justify-center py-4 transition-all duration-300",
                    isScrolled ? "py-4" : "py-6"
                )}
            >
                <div className={cn(
                    "w-[90%] max-w-7xl flex items-center justify-between px-6 py-3 rounded-full transition-all duration-300",
                    isScrolled
                        ? "bg-white/70 backdrop-blur-md shadow-lg border border-white/20"
                        : "bg-transparent"
                )}>
                    {/* Logo */}
                    <a href="#" className="text-2xl font-bold tracking-tight text-slate-800">
                        Schedule<span className="text-primary-600">Pro</span>
                    </a>

                    {/* Desktop Menu */}
                    <div className="hidden md:flex items-center gap-8">
                        {['Features', 'Integrations', 'Pricing', 'About'].map((item) => (
                            <a
                                key={item}
                                href={`#${item.toLowerCase()}`}
                                className="text-sm font-medium text-slate-600 hover:text-primary-600 transition-colors"
                            >
                                {item}
                            </a>
                        ))}
                    </div>

                    {/* Call to Action */}
                    <div className="hidden md:flex items-center gap-4">
                        <button className="text-sm font-semibold text-slate-700 hover:text-primary-600 transition-colors">
                            Log in
                        </button>
                        <button className="px-5 py-2.5 text-sm font-semibold text-white bg-primary-600 rounded-full hover:bg-primary-700 transition-all shadow-lg shadow-primary-500/30 hover:shadow-primary-500/50">
                            Get Started
                        </button>
                    </div>

                    {/* Mobile Menu Toggle */}
                    <button
                        className="md:hidden p-2 text-slate-600"
                        onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                    >
                        {isMobileMenuOpen ? <X /> : <Menu />}
                    </button>
                </div>
            </motion.nav>

            {/* Mobile Menu Overlay */}
            <AnimatePresence>
                {isMobileMenuOpen && (
                    <motion.div
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        className="fixed inset-0 z-40 bg-white/95 backdrop-blur-xl md:hidden pt-24 px-6"
                    >
                        <div className="flex flex-col gap-6 text-center">
                            {['Features', 'Integrations', 'Pricing', 'About'].map((item) => (
                                <a
                                    key={item}
                                    href={`#${item.toLowerCase()}`}
                                    onClick={() => setIsMobileMenuOpen(false)}
                                    className="text-2xl font-semibold text-slate-800 hover:text-primary-600"
                                >
                                    {item}
                                </a>
                            ))}
                            <div className="flex flex-col gap-4 mt-8">
                                <button className="w-full py-3 text-lg font-semibold text-slate-700 border border-slate-200 rounded-xl">
                                    Log in
                                </button>
                                <button className="w-full py-3 text-lg font-semibold text-white bg-primary-600 rounded-xl shadow-lg">
                                    Get Started
                                </button>
                            </div>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </>
    );
};

export default Navbar;
