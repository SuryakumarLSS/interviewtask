import { useEffect, useState } from 'react';
import { fetchRoles, createRole, fetchPermissions, savePermissions, fetchUsers, createUser, deletePermission, updateUserRole, deleteUser } from '../services/api';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowLeft, ChevronDown, Send, Trash2 } from 'lucide-react';
import { Checkbox } from '../components/ui';
import { toast } from 'sonner';
import { useAuth } from '../context/AuthContext';

export default function Admin() {
    const navigate = useNavigate();
    const { user } = useAuth();
    const [roles, setRoles] = useState([]);
    const [users, setUsers] = useState([]);

    const [selectedRole, setSelectedRole] = useState(null);
    const [permissions, setPermissions] = useState([]);
    const [newRoleName, setNewRoleName] = useState('');
    const [newUser, setNewUser] = useState({ username: '', role_id: 1 });
    const [deleteConfirmation, setDeleteConfirmation] = useState({ isOpen: false, userId: null, username: '' });

    const resources = ['employees', 'projects', 'orders'];
    const actions = ['read', 'create', 'update', 'delete'];

    useEffect(() => {
        loadRoles();
        loadUsers();
    }, []);

    useEffect(() => {
        if (selectedRole) {
            loadPermissions(selectedRole);
        }
    }, [selectedRole]);

    useEffect(() => {
        if (roles.length > 0) {
            const adminRole = roles.find(r => r.name.toLowerCase() === 'admin');
            if (adminRole) {
                setNewUser(prev => ({ ...prev, role_id: adminRole.id }));
            }
        }
    }, [roles]);

    const loadRoles = async () => {
        const data = await fetchRoles();
        setRoles(data);
    };

    const loadUsers = async () => {
        const data = await fetchUsers();
        setUsers(data);
    };

    const loadPermissions = async (roleId) => {
        const data = await fetchPermissions(roleId);
        setPermissions(data);
    };

    const handleCreateRole = async () => {
        if (!newRoleName) return;
        try {
            await createRole(newRoleName);
            setNewRoleName('');
            loadRoles();
            toast.success("Role created");
        } catch (e) {
            toast.error(e.response?.data?.error || "Error creating role");
        }
    };

    const handleCreateUser = async () => {
        if (!newUser.username) return;
        try {
            await createUser(newUser);
            const adminRole = roles.find(r => r.name.toLowerCase() === 'admin');
            setNewUser({ username: '', role_id: adminRole ? adminRole.id : 1 });
            loadUsers();
            toast.success("User created");
        } catch (e) {
            toast.error(e.response?.data?.error || "Error creating user");
        }
    };

    const handleDeleteUser = async () => {
        if (!deleteConfirmation.userId) return;
        try {
            await deleteUser(deleteConfirmation.userId);
            toast.success("User deleted");
            setDeleteConfirmation({ isOpen: false, userId: null, username: '' });
            loadUsers();
        } catch (err) {
            toast.error(err.response?.data?.message || "Failed to delete user");
        }
    };

    const savePermissionState = async (resource, action, enabled) => {
        if (!selectedRole) return;
        const toastId = toast.loading('Updating permission...');
        try {
            if (enabled) {
                await savePermissions({ role_id: selectedRole, resource, action });
            } else {
                await deletePermission({ role_id: selectedRole, resource, action });
            }
            loadPermissions(selectedRole);
            toast.success("Permission updated", { id: toastId });
        } catch (e) {
            toast.error(e.response?.data?.message || "Failed to update permission", { id: toastId });
        }
    };

    return (
        <div className="min-h-screen bg-zinc-50 font-sans text-zinc-900 pb-20 relative">
            <header className="px-8 py-6 bg-white border-b border-zinc-200 shadow-sm sticky top-0 z-30 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold tracking-tight text-zinc-900">Admin Console</h1>
                    <p className="text-zinc-500 text-sm mt-1">Manage system access, roles, and users</p>
                </div>
                <button
                    onClick={() => navigate(-1)}
                    className="flex items-center space-x-2 text-sm font-medium text-zinc-500 hover:text-zinc-900 transition-colors px-4 py-2 rounded-lg hover:bg-zinc-100"
                >
                    <ArrowLeft className="w-5 h-5" />
                    <span>Back</span>
                </button>
            </header>

            <main className="max-w-7xl mx-auto px-8 py-10 space-y-12">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.1 }}
                        className="bg-white rounded-xl shadow-[0_2px_10px_-1px_rgba(0,0,0,0.1)] border border-zinc-100 p-6"
                    >
                        <div className="flex items-center justify-between mb-6">
                            <h2 className="text-lg font-bold text-zinc-900">Roles</h2>
                            <div className="flex space-x-2">
                                <input
                                    value={newRoleName}
                                    onChange={(e) => setNewRoleName(e.target.value)}
                                    placeholder="New Role Name"
                                    className="border border-zinc-200 bg-zinc-50 rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:ring-1 focus:ring-zinc-900"
                                />
                                <button
                                    onClick={handleCreateRole}
                                    className="bg-zinc-900 text-white px-3 py-1.5 rounded-lg text-sm font-medium hover:bg-zinc-800 transition-colors"
                                >
                                    + Add
                                </button>
                            </div>
                        </div>
                        <div className="space-y-2 max-h-60 overflow-y-auto pr-2 custom-scrollbar">
                            {roles.map(role => (
                                <button
                                    key={role.id}
                                    onClick={() => setSelectedRole(role.id)}
                                    className={`w-full text-left px-4 py-3 rounded-lg flex justify-between items-center transition-all duration-200 border-l-4
                                        ${selectedRole === role.id
                                            ? 'bg-zinc-50 border-zinc-900 shadow-sm'
                                            : 'bg-white border-transparent hover:bg-zinc-50 hover:border-zinc-200'
                                        }
                                    `}
                                >
                                    <span className="font-medium text-sm text-zinc-900">{role.name}</span>
                                    <span className="text-xs text-zinc-400 font-mono">ID: {role.id}</span>
                                </button>
                            ))}
                        </div>
                    </motion.div>

                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.2 }}
                        className="bg-white rounded-xl shadow-[0_2px_10px_-1px_rgba(0,0,0,0.1)] border border-zinc-100 p-6"
                    >
                        <h2 className="text-lg font-bold text-zinc-900 mb-6">Invite User</h2>
                        <div className="space-y-4">
                            <input
                                value={newUser.username}
                                onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
                                placeholder="Email Address"
                                className="w-full px-4 py-3 bg-zinc-50 border border-zinc-200 rounded-lg focus:ring-2 focus:ring-zinc-900 focus:border-transparent outline-none transition-all text-sm font-medium"
                            />

                            <div className="relative">
                                <select
                                    value={newUser.role_id}
                                    onChange={(e) => setNewUser({ ...newUser, role_id: Number(e.target.value) })}
                                    className="w-full px-4 py-3 bg-zinc-50 border border-zinc-200 rounded-lg focus:ring-2 focus:ring-zinc-900 focus:border-transparent outline-none transition-all text-sm font-medium appearance-none cursor-pointer"
                                >
                                    <option value={0} disabled>Select Role</option>
                                    {roles.map(role => (
                                        <option key={role.id} value={role.id}>{role.name}</option>
                                    ))}
                                </select>
                                <div className="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-zinc-500">
                                    <ChevronDown className="w-4 h-4" />
                                </div>
                            </div>
                            <button
                                onClick={async () => {
                                    if (!newUser.username || !newUser.role_id) {
                                        toast.error("Please enter email and select a role");
                                        return;
                                    }
                                    const toastId = toast.loading('Sending invitation...');
                                    try {
                                        await createUser({
                                            email: newUser.username,
                                            role_id: newUser.role_id
                                        });
                                        const adminRole = roles.find(r => r.name.toLowerCase() === 'admin');
                                        setNewUser({ username: '', role_id: adminRole ? adminRole.id : 1 });
                                        loadUsers();
                                        toast.success("Invitation sent successfully", { id: toastId });
                                    } catch (e) {
                                        toast.error(e.response?.data?.message || "Failed to invite user", { id: toastId });
                                    }
                                }}
                                className="w-full px-6 py-3 bg-zinc-900 text-white rounded-lg text-sm font-bold hover:bg-zinc-800 transition-all shadow-lg shadow-zinc-900/10 active:scale-95 flex items-center justify-center"
                            >
                                <Send className="w-4 h-4 mr-2" />
                                Send Invitation
                            </button>
                        </div>

                        <div className="mt-8 pt-6 border-t border-zinc-100">
                            <h3 className="text-xs font-semibold text-zinc-400 uppercase tracking-wider mb-4">Current Users</h3>
                            <div className="max-h-60 overflow-y-auto pr-2 custom-scrollbar space-y-2">
                                {users.map(u => (
                                    <div key={u.id} className="flex items-center justify-between p-3 rounded-lg hover:bg-zinc-50 transition-colors border border-transparent hover:border-zinc-100">
                                        <div className="flex items-center space-x-3">
                                            <div className="w-8 h-8 rounded-full bg-zinc-200 flex items-center justify-center text-xs font-bold text-zinc-600">
                                                {u.username.slice(0, 1).toUpperCase()}
                                            </div>
                                            <div>
                                                <div className="text-sm font-medium text-zinc-800">{u.username}</div>
                                                <div className="text-xs text-zinc-400">{u.email || ''}</div>
                                            </div>
                                        </div>

                                        <div className="flex items-center space-x-2">
                                            {user?.username === 'admin' ? (
                                                u.username === 'admin' ? (
                                                    <span className="text-xs px-2 py-1 bg-zinc-100 text-zinc-600 rounded font-medium">
                                                        Admin
                                                    </span>
                                                ) : (
                                                    <div className="relative group">
                                                        <select
                                                            value={u.role_id}
                                                            onChange={async (e) => {
                                                                try {
                                                                    await updateUserRole(u.id, Number(e.target.value));
                                                                    toast.success("Role updated");
                                                                    loadUsers();
                                                                } catch (err) {
                                                                    toast.error(err.response?.data?.message || "Failed to update role");
                                                                }
                                                            }}
                                                            className="appearance-none bg-zinc-100 text-zinc-700 text-xs font-medium px-2 py-1 pr-6 rounded focus:outline-none focus:ring-1 focus:ring-zinc-300 cursor-pointer"
                                                        >
                                                            {roles.map(r => (
                                                                <option key={r.id} value={r.id}>{r.name}</option>
                                                            ))}
                                                        </select>
                                                        <ChevronDown className="w-3 h-3 text-zinc-500 absolute right-1 top-1.5 pointer-events-none" />
                                                    </div>
                                                )
                                            ) : (
                                                <span className="text-xs px-2 py-1 bg-zinc-100 text-zinc-600 rounded font-medium">
                                                    {roles.find(r => r.id === u.role_id)?.name || 'User'}
                                                </span>
                                            )}

                                            <span className={`text-[10px] px-1.5 py-0.5 rounded border ${u.status === 'Pending'
                                                ? 'bg-amber-50 text-amber-600 border-amber-100'
                                                : u.status === 'Declined'
                                                    ? 'bg-red-50 text-red-600 border-red-100'
                                                    : 'bg-emerald-50 text-emerald-600 border-emerald-100'
                                                }`}>
                                                {u.status || 'Active'}
                                            </span>

                                            {user?.username === 'admin' && u.username !== 'admin' && (
                                                <button
                                                    onClick={() => setDeleteConfirmation({ isOpen: true, userId: u.id, username: u.username })}
                                                    className="p-1.5 text-zinc-400 hover:text-red-500 hover:bg-red-50 rounded-md transition-colors"
                                                    title="Delete User"
                                                >
                                                    <Trash2 className="w-4 h-4" />
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </motion.div>
                </div>

                <AnimatePresence mode='wait'>
                    {selectedRole ? (
                        <motion.div
                            key={selectedRole}
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0, y: 20 }}
                            transition={{ duration: 0.3 }}
                            className="bg-white rounded-xl shadow-[0_2px_10px_-1px_rgba(0,0,0,0.1)] border border-zinc-100 overflow-hidden"
                        >
                            <div className="px-8 py-6 border-b border-zinc-100 flex items-center justify-between bg-zinc-50/50">
                                <h2 className="text-lg font-bold text-zinc-900">
                                    Permissions Configuration
                                    <span className="text-zinc-500 font-normal text-sm ml-2">for {roles.find(r => r.id === selectedRole)?.name}</span>
                                </h2>
                            </div>

                            <div className="overflow-x-auto">
                                <table className="w-full text-left">
                                    <thead className="bg-zinc-50/80 border-b border-zinc-100">
                                        <tr>
                                            <th className="px-8 py-4 text-xs font-semibold text-zinc-500 uppercase tracking-wider w-1/5">Resource</th>
                                            {actions.map(action => (
                                                <th key={action} className="px-6 py-4 text-xs font-semibold text-zinc-500 uppercase tracking-wider text-center w-1/5">
                                                    {action}
                                                </th>
                                            ))}
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-zinc-100">
                                        {resources.map(res => (
                                            <tr key={res} className="hover:bg-zinc-50/50 transition-colors">
                                                <td className="px-8 py-6">
                                                    <span className="font-bold text-zinc-800 capitalize text-sm">{res}</span>
                                                </td>
                                                {actions.map(action => {
                                                    const perm = permissions.find(p => p.resource === res && p.action === action);
                                                    return (
                                                        <td key={action} className="px-6 py-6 border-l border-zinc-50">
                                                            <div className="flex flex-col items-center space-y-3">
                                                                <Checkbox
                                                                    checked={!!perm}
                                                                    onChange={(checked) => savePermissionState(res, action, checked)}
                                                                    label={perm ? "Allowed" : "Restricted"}
                                                                />
                                                            </div>
                                                        </td>
                                                    );
                                                })}
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </motion.div>
                    ) : (
                        <div className="h-48 flex flex-col items-center justify-center text-zinc-400 border-2 border-dashed border-zinc-200 rounded-xl bg-zinc-50/50">
                            <span className="text-sm font-medium">Select a role above to configure permissions</span>
                        </div>
                    )}
                </AnimatePresence>
            </main>

            <AnimatePresence>
                {deleteConfirmation.isOpen && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            className="absolute inset-0 bg-black/20 backdrop-blur-sm"
                            onClick={() => setDeleteConfirmation({ isOpen: false, userId: null, username: '' })}
                        />
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            className="bg-white rounded-xl shadow-xl w-full max-w-sm relative z-10 overflow-hidden"
                        >
                            <div className="p-6 text-center">
                                <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                                    <Trash2 className="w-6 h-6 text-red-600" />
                                </div>
                                <h3 className="text-lg font-bold text-zinc-900 mb-2">Delete User?</h3>
                                <p className="text-sm text-zinc-500 mb-6">
                                    Are you sure you want to delete <span className="font-bold text-zinc-800">{deleteConfirmation.username}</span>? This action cannot be undone.
                                </p>
                                <div className="flex space-x-3">
                                    <button
                                        onClick={() => setDeleteConfirmation({ isOpen: false, userId: null, username: '' })}
                                        className="flex-1 py-2.5 bg-zinc-100 text-zinc-700 rounded-lg text-sm font-medium hover:bg-zinc-200 transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        onClick={handleDeleteUser}
                                        className="flex-1 py-2.5 bg-red-600 text-white rounded-lg text-sm font-medium hover:bg-red-700 transition-colors shadow-lg shadow-red-600/20"
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
