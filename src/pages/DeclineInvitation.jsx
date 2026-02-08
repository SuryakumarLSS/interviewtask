import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { declineInvitation } from '../services/api';
import { motion } from 'framer-motion';
import { CheckCircle, XCircle, AlertCircle } from 'lucide-react';

export default function DeclineInvitation() {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const token = searchParams.get('token');
    const [status, setStatus] = useState('loading');
    const [message, setMessage] = useState('');

    useEffect(() => {
        if (!token) {
            setStatus('error');
            setMessage('Invalid invitation link.');
            return;
        }

        const handleDecline = async () => {
            try {
                await declineInvitation(token);
                setStatus('success');
            } catch (err) {
                setStatus('error');
                setMessage(err.response?.data?.message || 'Failed to decline invitation.');
            }
        };

        handleDecline();
    }, [token]);

    return (
        <div className="min-h-screen flex items-center justify-center bg-zinc-50 font-sans p-4">
            <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="bg-white p-8 rounded-2xl shadow-xl w-full max-w-md text-center border border-zinc-100"
            >
                {status === 'loading' && (
                    <div className="flex flex-col items-center">
                        <div className="w-10 h-10 border-4 border-zinc-200 border-t-zinc-900 rounded-full animate-spin mb-4" />
                        <p className="text-zinc-600 font-medium">Processing your request...</p>
                    </div>
                )}

                {status === 'success' && (
                    <div className="flex flex-col items-center">
                        <div className="w-16 h-16 bg-green-50 rounded-full flex items-center justify-center mb-4">
                            <CheckCircle className="w-8 h-8 text-green-500" />
                        </div>
                        <h2 className="text-xl font-bold text-zinc-900 mb-2">Invitation Declined</h2>
                        <p className="text-zinc-500 mb-6">
                            You have successfully declined the invitation. No further action is required.
                        </p>
                        <button
                            onClick={() => navigate('/login')}
                            className="text-sm font-medium text-zinc-900 hover:text-zinc-700 underline"
                        >
                            Go to Login
                        </button>
                    </div>
                )}

                {status === 'error' && (
                    <div className="flex flex-col items-center">
                        <div className="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center mb-4">
                            {message.includes('Invalid') ? (
                                <AlertCircle className="w-8 h-8 text-red-500" />
                            ) : (
                                <XCircle className="w-8 h-8 text-red-500" />
                            )}
                        </div>
                        <h2 className="text-xl font-bold text-zinc-900 mb-2">Action Failed</h2>
                        <p className="text-zinc-500 mb-6">{message}</p>
                        <button
                            onClick={() => navigate('/login')}
                            className="bg-zinc-900 text-white px-6 py-2 rounded-lg text-sm font-medium hover:bg-zinc-800 transition-colors"
                        >
                            Go to Home
                        </button>
                    </div>
                )}
            </motion.div>
        </div>
    );
}
