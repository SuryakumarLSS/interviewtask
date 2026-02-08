import { motion } from 'framer-motion';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs) {
    return twMerge(clsx(inputs));
}

export const Checkbox = ({ className, checked, onChange, label }) => {
    return (
        <label className={cn("inline-flex items-center space-x-2 cursor-pointer group", className)}>
            <div className="relative">
                <input
                    type="checkbox"
                    className="appearance-none peer sr-only"
                    checked={checked}
                    onChange={(e) => onChange?.(e.target.checked)}
                />
                <motion.div
                    initial={false}
                    animate={{
                        backgroundColor: checked ? "#18181b" : "#fff",
                        borderColor: checked ? "#18181b" : "#e4e4e7"
                    }}
                    transition={{ duration: 0.2 }}
                    className="w-5 h-5 border-2 rounded-md flex items-center justify-center border-zinc-200 group-hover:border-zinc-400 peer-focus-visible:ring-2 peer-focus-visible:ring-offset-2 peer-focus-visible:ring-zinc-950"
                >
                    <motion.svg
                        initial={false}
                        animate={{
                            pathLength: checked ? 1 : 0,
                            opacity: checked ? 1 : 0
                        }}
                        transition={{ duration: 0.2, ease: "easeOut" }}
                        className="w-3.5 h-3.5 text-white stroke-current stroke-[3px]"
                        viewBox="0 0 24 24"
                        fill="none"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                    >
                        <path d="M20 6L9 17l-5-5" />
                    </motion.svg>
                </motion.div>
            </div>
            {label && <span className="text-sm font-medium text-zinc-700 select-none group-hover:text-zinc-900">{label}</span>}
        </label>
    );
};
