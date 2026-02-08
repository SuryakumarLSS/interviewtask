import { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { fetchMyPermissions, fetchResource, createResource, updateResource, deleteResource, fetchEmployeeNames } from '../services/api';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { LayoutDashboard, Trash2, AlertTriangle, X, Pencil, Plus, Users, Briefcase, Package, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';

export default function Dashboard() {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const [permissions, setPermissions] = useState([]);
    const [activeResource, setActiveResource] = useState(null);
    const [data, setData] = useState([]);
    const [headers, setHeaders] = useState([]);
    const [loading, setLoading] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);
    const [deleteModalOpen, setDeleteModalOpen] = useState(false);
    const [itemToDelete, setItemToDelete] = useState(null);
    const [currentItem, setCurrentItem] = useState(null);
    const [newItem, setNewItem] = useState({});

    const [employeesList, setEmployeesList] = useState([]);

    const canRead = (res) => user?.role_id === 1 || permissions.some(p => p.resource === res && p.action === 'read');
    const canCreate = (res) => user?.role_id === 1 || permissions.some(p => p.resource === res && p.action === 'create');
    const canUpdate = (res) => user?.role_id === 1 || permissions.some(p => p.resource === res && p.action === 'update');
    const canDelete = (res) => user?.role_id === 1 || permissions.some(p => p.resource === res && p.action === 'delete');

    const loadEmployees = async () => {
        try {
            let emps = await fetchEmployeeNames();
            if (!Array.isArray(emps)) {
                throw new Error("Invalid response format");
            }
            setEmployeesList(emps || []);
        } catch (err) {
            console.warn("Dedicated employee fetch failed, trying fallback...", err);
            try {
                const data = await fetchResource('employees');
                if (Array.isArray(data)) {
                    setEmployeesList(data);
                } else {
                    setEmployeesList([]);
                }
            } catch (fallbackErr) {
                console.error("All employee fetch methods failed", fallbackErr);
            }
        }
    };

    useEffect(() => {
        loadEmployees();
    }, []);

    useEffect(() => {
        if (modalOpen) {
            loadEmployees();
        }
    }, [modalOpen]);

    useEffect(() => {
        if (user?.role_id !== 1) {
            fetchMyPermissions().then(setPermissions).catch(console.error);
        } else {
            setPermissions([
                { resource: 'employees', action: 'read' },
                { resource: 'projects', action: 'read' },
                { resource: 'orders', action: 'read' }
            ]);
        }
    }, [user]);

    useEffect(() => {
        if (activeResource) {
            setLoading(true);
            fetchResource(activeResource)
                .then(resData => {
                    setData(resData);
                    if (resData.length > 0) {
                        setHeaders(Object.keys(resData[0]));
                    } else {
                        setHeaders([]);
                    }
                })
                .catch(err => console.error(err))
                .finally(() => setLoading(false));
        }
    }, [activeResource]);

    const RESOURCE_SCHEMAS = {
        employees: ['name', 'position', 'salary', 'department'],
        projects: ['name', 'assigned_to', 'status', 'budget'],
        orders: ['customer_name', 'amount', 'status', 'order_date']
    };

    const handleSave = async () => {
        if (!activeResource) return;
        const toastId = toast.loading('Saving...');
        try {
            const allowedFields = RESOURCE_SCHEMAS[activeResource] || headers;
            const filterData = (source) => {
                const filtered = {};
                allowedFields.forEach(field => {
                    if (source[field] !== undefined) filtered[field] = source[field];
                });
                return filtered;
            };

            if (currentItem) {
                const payload = filterData(currentItem);
                await updateResource(activeResource, currentItem.id, payload);
            } else {
                const payload = filterData(newItem);
                await createResource(activeResource, payload);
            }
            const updated = await fetchResource(activeResource);
            setData(updated);
            setModalOpen(false);
            setCurrentItem(null);
            setNewItem({});
            toast.success("Saved successfully!", { id: toastId });
        } catch (err) {
            toast.error(err.response?.data?.message || "Operation failed! Check permissions.", { id: toastId });
        }
    };

    const handleDeleteClick = (id) => {
        setItemToDelete(id);
        setDeleteModalOpen(true);
    };

    const confirmDelete = async () => {
        if (!activeResource || itemToDelete === null) return;

        const toastId = toast.loading('Deleting...');
        try {
            await deleteResource(activeResource, itemToDelete);
            setData(data.filter(item => item.id !== itemToDelete));
            toast.success("Deleted successfully!", { id: toastId });
            setDeleteModalOpen(false);
            setItemToDelete(null);
        } catch (err) {
            toast.error("Delete failed! Check permissions.", { id: toastId });
        }
    };

    const resources = ['employees', 'projects', 'orders'];
    const visibleResources = resources.filter(res => canRead(res));

    useEffect(() => {
        if (!activeResource && visibleResources.length > 0) {
            setActiveResource(visibleResources[0]);
        }
    }, [visibleResources, activeResource]);

    return (
        <div className="flex h-screen bg-zinc-50 font-sans text-zinc-900 overflow-hidden">
            <aside className="w-64 bg-white border-r border-zinc-200 flex flex-col z-20">
                <div className="h-16 flex items-center px-6 border-b border-zinc-100">
                    <div className="w-6 h-6 bg-zinc-900 rounded-md mr-3 flex items-center justify-center text-white text-xs font-bold">R</div>
                    <span className="font-bold tracking-tight text-lg">RBAC</span>
                </div>

                <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
                    <p className="px-3 text-xs font-semibold text-zinc-400 uppercase tracking-wider mb-2">Workspace</p>
                    {visibleResources.map(res => {
                        const iconMap = {
                            employees: Users,
                            projects: Briefcase,
                            orders: Package
                        };
                        const Icon = iconMap[res] || LayoutDashboard;

                        return (
                            <button
                                key={res}
                                onClick={() => setActiveResource(res)}
                                className={`w-full text-left px-3 py-2.5 rounded-md text-sm font-medium transition-all duration-200 flex items-center
                                    ${activeResource === res
                                        ? 'bg-zinc-100 text-zinc-900 shadow-sm'
                                        : 'text-zinc-500 hover:bg-zinc-50 hover:text-zinc-900'
                                    }
                                `}
                            >
                                <Icon className="w-4 h-4 mr-3" />
                                {res.charAt(0).toUpperCase() + res.slice(1)}
                            </button>
                        );
                    })}

                    <div className="pt-8">
                        {user?.role_id === 1 && (
                            <>
                                <p className="px-3 text-xs font-semibold text-zinc-400 uppercase tracking-wider mb-2">Settings</p>
                                <button
                                    onClick={() => navigate('/admin')}
                                    className="w-full text-left px-3 py-2.5 rounded-md text-sm font-medium text-zinc-500 hover:bg-zinc-50 hover:text-zinc-900 transition-all flex items-center"
                                >
                                    <ShieldCheck className="w-4 h-4 mr-3" />
                                    Admin Panel
                                </button>
                            </>
                        )}
                    </div>
                </nav>

                <div className="p-4 border-t border-zinc-100">
                    <div className="flex items-center space-x-3 mb-4">
                        <div className="w-8 h-8 rounded-full bg-zinc-100 flex items-center justify-center text-xs font-bold border border-zinc-200 uppercase">
                            {user?.username.slice(0, 2)}
                        </div>
                        <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium text-zinc-900 truncate">{user?.username}</p>
                            <p className="text-xs text-zinc-500 truncate">role_id: {user?.role_id}</p>
                        </div>
                    </div>
                    <button
                        onClick={logout}
                        className="w-full px-3 py-2 text-xs font-medium text-zinc-600 border border-zinc-200 rounded-md hover:bg-zinc-50 hover:text-red-600 transition-colors"
                    >
                        Log out
                    </button>
                </div>
            </aside>

            <main className="flex-1 flex flex-col min-w-0 overflow-hidden bg-zinc-50/50">
                <header className="h-16 px-8 flex items-center justify-between border-b border-zinc-200 bg-white/80 backdrop-blur-sm z-10">
                    <h1 className="text-xl font-bold tracking-tight text-zinc-900 capitalize">
                        {activeResource ? `${activeResource}` : 'Dashboard'}
                    </h1>
                    {activeResource && canCreate(activeResource) && (
                        <button
                            onClick={() => { setCurrentItem(null); setModalOpen(true); }}
                            className="bg-zinc-900 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-zinc-800 transition-all shadow-lg shadow-zinc-900/10 active:scale-95 flex items-center space-x-2"
                        >
                            <Plus className="w-4 h-4" />
                            <span>New Entry</span>
                        </button>
                    )}
                </header>

                <div className="flex-1 overflow-auto p-8 relative">
                    <AnimatePresence mode='wait'>
                        {activeResource ? (
                            <motion.div
                                key={activeResource}
                                initial={{ opacity: 0, scale: 0.98 }}
                                animate={{ opacity: 1, scale: 1 }}
                                exit={{ opacity: 0, scale: 0.98 }}
                                transition={{ duration: 0.2 }}
                                className="bg-white rounded-xl shadow-[0_2px_20px_-5px_rgba(0,0,0,0.1)] border border-zinc-100 overflow-hidden"
                            >
                                <div className="overflow-x-auto">
                                    <table className="w-full text-left">
                                        <thead className="bg-zinc-50/50 border-b border-zinc-100">
                                            <tr>
                                                {headers.map(h => (
                                                    <th key={h} className="px-6 py-4 font-semibold text-zinc-500 uppercase text-[10px] tracking-wider">
                                                        {h.replaceAll('_', ' ')}
                                                    </th>
                                                ))}
                                                {(canUpdate(activeResource) || canDelete(activeResource)) && (
                                                    <th className="px-6 py-4 font-semibold text-zinc-500 uppercase text-[10px] tracking-wider text-right">
                                                        Actions
                                                    </th>
                                                )}
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-zinc-50">
                                            {data.map((row, idx) => (
                                                <motion.tr
                                                    key={row.id || idx}
                                                    initial={{ opacity: 0 }}
                                                    animate={{ opacity: 1 }}
                                                    transition={{ delay: idx * 0.05 }}
                                                    className="group hover:bg-zinc-50/80 transition-colors"
                                                >
                                                    {headers.map(h => (
                                                        <td key={h} className="px-6 py-4 text-sm text-zinc-700 font-medium group-hover:text-zinc-900">
                                                            {row[h]}
                                                        </td>
                                                    ))}
                                                    {(canUpdate(activeResource) || canDelete(activeResource)) && (
                                                        <td className="px-6 py-4 text-right space-x-2">
                                                            {canUpdate(activeResource) && (
                                                                <button
                                                                    onClick={() => { setCurrentItem(row); setModalOpen(true); }}
                                                                    className="inline-flex items-center justify-center w-8 h-8 rounded-full text-zinc-400 hover:text-zinc-900 hover:bg-zinc-100 transition-all"
                                                                    title="Edit"
                                                                >
                                                                    <Pencil className="w-4 h-4" />
                                                                </button>
                                                            )}
                                                            {canDelete(activeResource) && (
                                                                <button
                                                                    onClick={() => handleDeleteClick(row.id)}
                                                                    className="inline-flex items-center justify-center w-8 h-8 rounded-full text-zinc-400 hover:text-red-600 hover:bg-red-50 transition-all"
                                                                    title="Delete"
                                                                >
                                                                    <Trash2 className="w-4 h-4" />
                                                                </button>
                                                            )}
                                                        </td>
                                                    )}
                                                </motion.tr>
                                            ))}
                                        </tbody>
                                    </table>
                                    {data.length === 0 && !loading && (
                                        <div className="p-12 text-center">
                                            <div className="w-12 h-12 bg-zinc-50 rounded-full flex items-center justify-center mx-auto mb-3">
                                                <LayoutDashboard className="w-6 h-6 text-zinc-300" />
                                            </div>
                                            <p className="text-zinc-400 text-sm">No records found.</p>
                                        </div>
                                    )}
                                </div>
                            </motion.div>
                        ) : (
                            <motion.div
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                className="flex flex-col items-center justify-center h-full text-zinc-400"
                            >
                                <div className="w-16 h-16 bg-zinc-100 rounded-2xl mb-4 flex items-center justify-center">
                                    <LayoutDashboard className="w-8 h-8 text-zinc-300" />
                                </div>
                                <p className="text-sm font-medium">Select a resource to manage</p>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            </main>

            <AnimatePresence>
                {modalOpen && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            className="absolute inset-0 bg-zinc-900/20 backdrop-blur-sm"
                            onClick={() => setModalOpen(false)}
                        />
                        <motion.div
                            initial={{ scale: 0.95, opacity: 0, y: 10 }}
                            animate={{ scale: 1, opacity: 1, y: 0 }}
                            exit={{ scale: 0.95, opacity: 0, y: 10 }}
                            className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-lg z-10 relative border border-zinc-100"
                        >
                            <button
                                onClick={() => setModalOpen(false)}
                                className="absolute top-4 right-4 text-zinc-400 hover:text-zinc-600 p-1"
                            >
                                <X className="w-5 h-5" />
                            </button>
                            <form onSubmit={(e) => { e.preventDefault(); handleSave(); }}>
                                <h3 className="text-xl font-bold tracking-tight mb-6 text-zinc-900">
                                    {currentItem ? `Edit entry` : `New entry`}
                                </h3>

                                <div className="space-y-5 max-h-[60vh] overflow-y-auto px-1 py-1 pr-2 custom-scrollbar">
                                    {(() => {
                                        const RESOURCE_SCHEMAS = {
                                            employees: ['name', 'position', 'salary', 'department'],
                                            projects: ['name', 'assigned_to', 'status', 'budget'],
                                            orders: ['customer_name', 'amount', 'status', 'order_date']
                                        };

                                        const fields = activeResource && RESOURCE_SCHEMAS[activeResource]
                                            ? RESOURCE_SCHEMAS[activeResource]
                                            : headers;

                                        return fields.map(field => {
                                            if (field === 'id') return null;

                                            if (field === 'assigned_to') {
                                                return (
                                                    <div key={field} className="relative">
                                                        <label className="block text-xs font-semibold text-zinc-500 uppercase tracking-wider mb-2">
                                                            {field.replace('_', ' ')} <span className="text-red-500">*</span>
                                                        </label>
                                                        <div className="relative">
                                                            <select
                                                                required
                                                                value={currentItem ? currentItem[field] || '' : newItem[field] || ''}
                                                                onChange={(e) => {
                                                                    const val = e.target.value;
                                                                    if (currentItem) setCurrentItem({ ...currentItem, [field]: val });
                                                                    else setNewItem({ ...newItem, [field]: val });
                                                                }}
                                                                className="w-full px-4 py-2.5 bg-zinc-50 border border-zinc-200 rounded-lg focus:ring-2 focus:ring-zinc-900 focus:border-transparent outline-none transition-all text-sm font-medium text-zinc-900 appearance-none cursor-pointer"
                                                            >
                                                                <option value="" disabled>Select employee...</option>
                                                                {employeesList.length > 0 ? (
                                                                    employeesList.map(u => (
                                                                        <option key={u.id} value={u.name}>
                                                                            {u.name}
                                                                        </option>
                                                                    ))
                                                                ) : (
                                                                    <option value="" disabled>No employees found (Create one first)</option>
                                                                )}
                                                            </select>
                                                            <div className="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-zinc-500">
                                                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                                                                </svg>
                                                            </div>
                                                        </div>
                                                    </div>
                                                );
                                            }

                                            if (field === 'status') {
                                                const currentVal = (currentItem ? currentItem[field] : newItem[field]) || '';

                                                const getStatusColor = (val) => {
                                                    switch (val.toLowerCase()) {
                                                        case 'open': return 'bg-blue-50 text-blue-700 border-blue-200';
                                                        case 'in progress': return 'bg-amber-50 text-amber-700 border-amber-200';
                                                        case 'pending': return 'bg-zinc-50 text-zinc-700 border-zinc-200';
                                                        default: return 'bg-zinc-50 border-zinc-200 text-zinc-900';
                                                    }
                                                };

                                                return (
                                                    <div key={field} className="relative">
                                                        <label className="block text-xs font-semibold text-zinc-500 uppercase tracking-wider mb-2">
                                                            {field.replace('_', ' ')} <span className="text-red-500">*</span>
                                                        </label>
                                                        <div className="relative">
                                                            <select
                                                                required
                                                                value={currentVal}
                                                                onChange={(e) => {
                                                                    const val = e.target.value;
                                                                    if (currentItem) setCurrentItem({ ...currentItem, [field]: val });
                                                                    else setNewItem({ ...newItem, [field]: val });
                                                                }}
                                                                className={`w-full px-4 py-2.5 rounded-lg border focus:ring-2 focus:ring-zinc-900 focus:border-transparent outline-none transition-all text-sm font-semibold appearance-none cursor-pointer ${getStatusColor(currentVal)}`}
                                                            >
                                                                <option value="" disabled>Select status...</option>
                                                                <option value="Open">🔵 &nbsp; Open</option>
                                                                <option value="In Progress">🔄 &nbsp; In Progress</option>
                                                                <option value="Pending">⏳ &nbsp; Pending</option>
                                                            </select>
                                                            <div className="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none opacity-50">
                                                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                                                                </svg>
                                                            </div>
                                                        </div>
                                                    </div>
                                                );
                                            }

                                            return (
                                                <div key={field}>
                                                    <label className="block text-xs font-semibold text-zinc-500 uppercase tracking-wider mb-2">
                                                        {field.replace('_', ' ')} <span className="text-red-500">*</span>
                                                    </label>
                                                    <input
                                                        required
                                                        type={['salary', 'amount', 'budget'].includes(field) ? 'number' : 'text'}
                                                        defaultValue={currentItem ? currentItem[field] : ''}
                                                        onChange={(e) => {
                                                            if (currentItem) setCurrentItem({ ...currentItem, [field]: e.target.value });
                                                            else setNewItem({ ...newItem, [field]: e.target.value });
                                                        }}
                                                        className="w-full px-4 py-2.5 bg-zinc-50 border border-zinc-200 rounded-lg focus:ring-2 focus:ring-zinc-900 focus:border-transparent outline-none transition-all text-sm font-medium text-zinc-900 placeholder-zinc-400"
                                                        placeholder={`Enter ${field.replace('_', ' ')}`}
                                                    />
                                                </div>
                                            );
                                        });
                                    })()}
                                </div>

                                <div className="mt-8 flex justify-end space-x-3 pt-4 border-t border-zinc-100">
                                    <button
                                        type="button"
                                        onClick={() => setModalOpen(false)}
                                        className="px-4 py-2 text-sm font-medium text-zinc-500 hover:text-zinc-900 transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type="submit"
                                        className="px-6 py-2 bg-zinc-900 text-white rounded-lg text-sm font-medium hover:bg-zinc-800 transition-all shadow-lg shadow-zinc-900/10 active:scale-95"
                                    >
                                        Save Changes
                                    </button>
                                </div>
                            </form>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>

            <AnimatePresence>
                {deleteModalOpen && (
                    <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            className="absolute inset-0 bg-zinc-900/20 backdrop-blur-sm"
                            onClick={() => setDeleteModalOpen(false)}
                        />
                        <motion.div
                            initial={{ scale: 0.95, opacity: 0 }}
                            animate={{ scale: 1, opacity: 1 }}
                            exit={{ scale: 0.95, opacity: 0 }}
                            className="bg-white rounded-2xl shadow-xl w-full max-w-sm z-10 relative overflow-hidden"
                        >
                            <div className="p-6 text-center">
                                <div className="w-12 h-12 bg-red-100/50 rounded-full flex items-center justify-center mx-auto mb-4">
                                    <AlertTriangle className="w-6 h-6 text-red-600" />
                                </div>
                                <h3 className="text-lg font-bold text-zinc-900 mb-2">Delete Entry?</h3>
                                <p className="text-sm text-zinc-500 mb-6">
                                    Are you sure you want to delete this item? This action cannot be undone.
                                </p>
                                <div className="flex space-x-3">
                                    <button
                                        onClick={() => setDeleteModalOpen(false)}
                                        className="flex-1 px-4 py-2 bg-zinc-100 text-zinc-700 rounded-lg text-sm font-medium hover:bg-zinc-200 transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        onClick={confirmDelete}
                                        className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg text-sm font-medium hover:bg-red-700 transition-colors shadow-lg shadow-red-900/20"
                                    >
                                        Delete
                                    </button>
                                </div>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </div>
    );
}
