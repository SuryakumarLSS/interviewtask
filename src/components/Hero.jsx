import React from 'react';
import { motion } from 'framer-motion';
import { Calendar, Clock, ArrowRight, Zap, Sparkles, Command } from 'lucide-react';
import CosmicDust from './CosmicDust';

const Hero = () => {
    return (
        <section className="relative min-h-screen pt-32 pb-20 lg:pt-48 lg:pb-32 overflow-hidden bg-slate-50 selection:bg-purple-100 selection:text-purple-900">
            {/* Cosmic Background Effects */}
            <CosmicDust />

            {/* Gradient Orbs */}
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full z-0 pointer-events-none">
                <div className="absolute top-[10%] left-[20%] w-[600px] h-[600px] bg-indigo-200/40 rounded-full blur-[120px] mix-blend-multiply animate-float" />
                <div className="absolute top-[20%] right-[20%] w-[500px] h-[500px] bg-purple-200/40 rounded-full blur-[100px] mix-blend-multiply animation-delay-2000 animate-float-delayed" />
            </div>

            <div className="container relative z-10 mx-auto px-6">
                <div className="flex flex-col items-center text-center max-w-5xl mx-auto">

                    <motion.div
                        initial={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.5 }}
                        className="inline-flex items-center gap-2 px-6 py-2 rounded-full bg-white/80 backdrop-blur-md border border-purple-100 shadow-lg shadow-purple-500/5 mb-8 hover:scale-105 transition-transform cursor-default"
                    >
                        <Sparkles className="w-4 h-4 text-purple-600 animate-pulse" />
                        <span className="text-sm font-semibold bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-indigo-600">
                            The Future of Scheduling is Here
                        </span>
                    </motion.div>

                    <motion.h1
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6, delay: 0.1 }}
                        className="text-6xl md:text-8xl font-black tracking-tighter text-slate-900 mb-8 leading-none"
                    >
                        Time is <br className="md:hidden" />
                        <span className="relative inline-block">
                            <span className="absolute -inset-1 bg-gradient-to-r from-purple-600 via-pink-500 to-indigo-600 blur opacity-30 animate-pulse-glow" />
                            <span className="relative bg-clip-text text-transparent bg-gradient-to-r from-purple-600 via-pink-500 to-indigo-600">
                                Infinite
                            </span>
                        </span>
                    </motion.h1>

                    <motion.p
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6, delay: 0.2 }}
                        className="text-xl md:text-2xl text-slate-600 mb-12 max-w-2xl leading-relaxed"
                    >
                        Experience the fluid harmony of AI-driven scheduling.
                        No boxes. No conflicts. Just <span className="text-slate-900 font-semibold">pure flow</span>.
                    </motion.p>

                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6, delay: 0.3 }}
                        className="flex flex-col sm:flex-row items-center gap-5"
                    >
                        <button className="px-8 py-4 bg-slate-900 text-white font-bold rounded-full hover:bg-slate-800 transition-all shadow-xl hover:shadow-2xl hover:-translate-y-1 flex items-center gap-2 group">
                            Start Free Trial
                            <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                        </button>
                        <button className="px-8 py-4 bg-white/50 backdrop-blur-sm text-slate-700 font-bold rounded-full border border-white/60 hover:bg-white hover:border-white transition-all shadow-lg hover:shadow-xl flex items-center gap-2">
                            <Command className="w-5 h-5" />
                            Watch the Film
                        </button>
                    </motion.div>
                </div>

                {/* Fluid Abstract Mocks */}
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 1, delay: 0.5 }}
                    className="mt-24 relative h-[500px] w-full max-w-6xl mx-auto"
                >
                    {/* Central Floating Glass Panel */}
                    <motion.div
                        animate={{ y: [0, -20, 0] }}
                        transition={{ duration: 6, repeat: Infinity, ease: "easeInOut" }}
                        className="absolute left-1/2 -translate-x-1/2 top-0 z-20"
                    >
                        <div className="w-[320px] md:w-[600px] bg-white/70 backdrop-blur-xl border border-white/40 rounded-3xl p-6 shadow-2xl ring-1 ring-purple-500/10">
                            <div className="flex items-center justify-between mb-6">
                                <div className="flex items-center gap-3">
                                    <div className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-500 to-indigo-600 flex items-center justify-center text-white">
                                        <Zap className="w-5 h-5" />
                                    </div>
                                    <div>
                                        <h3 className="font-bold text-slate-800">Product Review</h3>
                                        <p className="text-xs text-slate-500 font-medium tracking-wide">AI NEGOTIATING...</p>
                                    </div>
                                </div>
                                <span className="px-3 py-1 rounded-full bg-green-100 text-green-700 text-xs font-bold uppercase tracking-wider">Confirmed</span>
                            </div>
                            <div className="space-y-3">
                                <div className="h-2 w-3/4 bg-slate-200/50 rounded-full" />
                                <div className="h-2 w-1/2 bg-slate-200/50 rounded-full" />
                            </div>
                            <div className="mt-6 flex items-center gap-2">
                                <div className="flex -space-x-2">
                                    {[1, 2, 3].map(i => (
                                        <div key={i} className="w-8 h-8 rounded-full border-2 border-white bg-slate-200 flex items-center justify-center text-[10px] text-slate-500 font-bold">U{i}</div>
                                    ))}
                                </div>
                                <div className="text-xs text-slate-400 font-medium">+ 2 others joined</div>
                            </div>
                        </div>
                    </motion.div>

                    {/* Floating Elements - Not Boxy Cards */}
                    <motion.div
                        animate={{ y: [0, 30, 0], x: [0, 10, 0] }}
                        transition={{ duration: 8, repeat: Infinity, ease: "easeInOut", delay: 1 }}
                        className="absolute left-[5%] md:left-[10%] top-[20%] z-10"
                    >
                        <div className="bg-white/40 backdrop-blur-lg border border-white/30 p-4 rounded-2xl shadow-xl">
                            <Calendar className="w-8 h-8 text-pink-500 mb-2" />
                            <div className="text-2xl font-bold text-slate-800">Oct 24</div>
                        </div>
                    </motion.div>

                    <motion.div
                        animate={{ y: [0, -25, 0], x: [0, -10, 0] }}
                        transition={{ duration: 7, repeat: Infinity, ease: "easeInOut", delay: 2 }}
                        className="absolute right-[5%] md:right-[15%] top-[40%] z-30"
                    >
                        <div className="bg-slate-900/90 backdrop-blur-md text-white p-5 rounded-2xl shadow-2xl flex items-center gap-4">
                            <Clock className="w-6 h-6 text-purple-400" />
                            <div>
                                <div className="font-bold">10:00 AM</div>
                                <div className="text-xs text-slate-400">PST (Local)</div>
                            </div>
                        </div>
                    </motion.div>

                </motion.div>
            </div>
        </section>
    );
};

export default Hero;
