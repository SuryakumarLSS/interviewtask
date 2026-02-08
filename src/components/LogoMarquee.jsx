import React from 'react';

const LogoMarquee = () => {
    // Using text for logos to keep it simple but styled professionally
    const logos = [
        { name: "Salesforce", color: "#00A1E0" },
        { name: "HubSpot", color: "#FF7A59" },
        { name: "Zoom", color: "#2D8CFF" },
        { name: "Slack", color: "#4A154B" },
        { name: "Microsoft", color: "#00A4EF" },
        { name: "Google", color: "#4285F4" },
        { name: "Stripe", color: "#008CDD" },
        { name: "Intercom", color: "#1F8CEB" },
        { name: "Notion", color: "#000000" },
        { name: "Zapier", color: "#FF4F00" },
        { name: "Linear", color: "#5E6AD2" },
        { name: "Calendly", color: "#006BFF" }
    ];

    return (
        <div className="w-full relative overflow-hidden py-10 bg-white/50 backdrop-blur-sm border-y border-slate-100/50">
            <div className="absolute inset-y-0 left-0 w-32 bg-gradient-to-r from-slate-50 to-transparent z-10" />
            <div className="absolute inset-y-0 right-0 w-32 bg-gradient-to-l from-slate-50 to-transparent z-10" />

            <div className="flex w-max animate-scroll">
                {[...logos, ...logos, ...logos].map((logo, idx) => (
                    <div
                        key={idx}
                        className="mx-8 flex items-center gap-3 opacity-60 hover:opacity-100 transition-opacity cursor-pointer group"
                    >
                        <span
                            className="text-2xl font-bold tracking-tight"
                            style={{ color: 'var(--slate-700)' }}
                        >
                            <span className="bg-clip-text text-transparent bg-gradient-to-tr from-slate-600 to-slate-400 group-hover:from-primary-600 group-hover:to-indigo-600 transition-all duration-300">
                                {logo.name}
                            </span>
                        </span>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default LogoMarquee;
