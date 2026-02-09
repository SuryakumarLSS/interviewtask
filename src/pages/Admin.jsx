import { useEffect, useState } from 'react';
import { fetchRoles, createRole, deleteRole, fetchPermissions, savePermissions, fetchUsers, createUser, deletePermission, updateUserRole, deleteUser, fetchRoleFieldPermissions, updateRoleFieldPermission } from '../services/api';
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
    const [fieldPermissions, setFieldPermissions] = useState([]);
    const [newRoleName, setNewRoleName] = useState('');
    const [newUser, setNewUser] = useState({ username: '', role_id: 1 });
    const [deleteConfirmation, setDeleteConfirmation] = useState({ isOpen: false, userId: null, username: '', type: 'user' });

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
        const [permsData, fieldPermsData] = await Promise.all([
            fetchPermissions(roleId),
            fetchRoleFieldPermissions(roleId)
        ]);
        setPermissions(permsData || []);
        setFieldPermissions(fieldPermsData || []);
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

    const handleDelete = async () => {
        if (!deleteConfirmation.userId) return;
        const isUser = deleteConfirmation.type === 'user';
        try {
            if (isUser) {
                await deleteUser(deleteConfirmation.userId);
                toast.success("User deleted");
            } else {
                await deleteRole(deleteConfirmation.userId);
                toast.success("Role deleted");
                if (selectedRole === deleteConfirmation.userId) {
                    setSelectedRole(null);
                }
            }
            setDeleteConfirmation({ isOpen: false, userId: null, username: '', type: 'user' });
            isUser ? loadUsers() : loadRoles();
        } catch (err) {
            toast.error(err.response?.data?.error || `Failed to delete ${deleteConfirmation.type}`);
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

    const handleFieldPermissionChange = async (resource, field, type, value) => {
        if (!selectedRole) return;

        // Optimistic UI update
        const updatedPerms = fieldPermissions.map(p => {
            if (p.resource === resource && p.field === field) {
                return { ...p, [type === 'view' ? 'can_view' : 'can_edit']: value };
            }
            return p;
        });
        setFieldPermissions(updatedPerms);

        const record = updatedPerms.find(p => p.resource === resource && p.field === field);

        try {
            await updateRoleFieldPermission({
                role_id: selectedRole,
                resource,
                field,
                can_view: record.can_view,
                can_edit: record.can_edit
            });
        } catch (e) {
            toast.error("Failed to update field permission");
            loadPermissions(selectedRole);
        }
    };

    const activeRoleName = roles.find(r => r.id === selectedRole)?.name;

    return (
        <div className="flex flex-col h-screen bg-slate-50 text-slate-900 font-sans overflow-hidden">
            {/* Top Navigation Bar */}
            <header className="flex-none h-16 bg-white border-b border-slate-200 px-6 flex items-center justify-between z-20 shadow-sm">
                <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-slate-900 rounded-lg flex items-center justify-center text-white font-bold">A</div>
                    <div>
                        <h1 className="text-lg font-bold text-slate-900 leading-tight">Admin Console</h1>
                        <p className="text-xs text-slate-500">Security & Access Control</p>
                    </div>
                </div>
                <button onClick={() => navigate(-1)} className="text-sm font-medium text-slate-500 hover:text-slate-900 flex items-center gap-2 px-3 py-2 rounded hover:bg-slate-100 transition-colors">
                    <ArrowLeft className="w-4 h-4" /> Back to App
                </button>
            </header>

            <div className="flex-1 flex overflow-hidden">
                {/* LEFT SIDEBAR: ROLES */}
                <aside className="w-64 bg-white border-r border-slate-200 flex flex-col z-10">
                    <div className="p-4 border-b border-slate-100">
                        <h2 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3">Roles</h2>
                        <div className="flex gap-2">
                            <input
                                value={newRoleName}
                                onChange={(e) => setNewRoleName(e.target.value)}
                                placeholder="New Role..."
                                className="flex-1 bg-slate-50 border border-slate-200 rounded text-sm px-2 py-1.5 focus:outline-none focus:border-slate-400"
                            />
                            <button onClick={handleCreateRole} className="bg-slate-900 text-white px-3 rounded text-sm font-medium hover:bg-slate-800">+</button>
                        </div>
                    </div>
                    <div className="flex-1 overflow-y-auto py-2">
                        {roles.map(role => (
                            <div
                                key={role.id}
                                className={`group flex items-center justify-between border-l-2 transition-all ${selectedRole === role.id ? 'bg-slate-50 border-slate-900' : 'border-transparent hover:bg-slate-50'
                                    }`}
                            >
                                <button
                                    onClick={() => setSelectedRole(role.id)}
                                    className={`flex-1 text-left px-4 py-3 ${selectedRole === role.id ? 'text-slate-900 font-medium' : 'text-slate-600 group-hover:text-slate-900'
                                        }`}
                                >
                                    <span className="text-sm">{role.name}</span>
                                </button>
                                {role.name.toLowerCase() !== 'admin' && (
                                    <button
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            setDeleteConfirmation({ isOpen: true, userId: role.id, username: role.name, type: 'role' });
                                        }}
                                        className="opacity-0 group-hover:opacity-100 p-2 mr-1 text-slate-400 hover:text-red-600 transition-opacity"
                                        title="Delete Role"
                                    >
                                        <Trash2 className="w-3.5 h-3.5" />
                                    </button>
                                )}
                                {selectedRole === role.id && role.name.toLowerCase() === 'admin' && <div className="w-1.5 h-1.5 rounded-full bg-slate-900 mr-4" />}
                            </div>
                        ))}
                    </div>
                </aside>

                {/* MAIN CONTENT AREA */}
                <main className="flex-1 flex flex-col bg-slate-50/50 overflow-hidden">
                    {selectedRole ? (
                        <div className="flex-1 overflow-y-auto p-8">
                            <div className="max-w-5xl mx-auto space-y-8">

                                {/* Header */}
                                <div className="flex items-center justify-between">
                                    <div>
                                        <h2 className="text-2xl font-bold text-slate-900">{activeRoleName}</h2>
                                        <p className="text-sm text-slate-500">Configure global and granular permissions for this role.</p>
                                    </div>
                                    <span className="px-3 py-1 bg-slate-200 text-slate-700 text-xs font-mono rounded">ID: {selectedRole}</span>
                                </div>

                                {/* 1. Resource Access (Matrix) */}
                                <section className="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
                                    <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/30">
                                        <div>
                                            <h3 className="text-sm font-bold text-slate-900">Table Level Permission</h3>
                                            <p className="text-xs text-slate-500">Global Resource Access</p>
                                        </div>
                                    </div>
                                    <table className="w-full text-left bg-white">
                                        <thead className="bg-slate-50 border-b border-slate-200/60">
                                            <tr>
                                                <th className="px-6 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider w-1/4">Resource</th>
                                                {actions.map(a => <th key={a} className="px-4 py-3 text-xs font-semibold text-slate-500 uppercase tracking-wider text-center flex-1">{a}</th>)}
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-slate-100">
                                            {resources.map(res => (
                                                <tr key={res} className="hover:bg-slate-50 transition-colors">
                                                    <td className="px-6 py-4 font-medium text-sm text-slate-700 capitalize">{res}</td>
                                                    {actions.map(action => {
                                                        const perm = (permissions || []).find(p => p.resource === res && p.action === action);
                                                        return (
                                                            <td key={action} className="px-4 py-4 text-center">
                                                                <label className="inline-flex items-center cursor-pointer">
                                                                    <input
                                                                        type="checkbox"
                                                                        className="peer sr-only"
                                                                        checked={!!perm}
                                                                        onChange={(e) => savePermissionState(res, action, e.target.checked)}
                                                                    />
                                                                    <div className={`
                                                                        w-9 h-5 rounded-full peer peer-focus:ring-2 peer-focus:ring-slate-300 
                                                                        peer-checked:after:translate-x-full peer-checked:after:border-white 
                                                                        after:content-[''] after:absolute after:top-[2px] after:left-[2px] 
                                                                        after:bg-white after:border-gray-300 after:border after:rounded-full 
                                                                        after:h-4 after:w-4 after:transition-all border-gray-200
                                                                        ${perm ? 'bg-slate-900 border-slate-900' : 'bg-gray-200'}
                                                                    `}></div>
                                                                </label>
                                                            </td>
                                                        );
                                                    })}
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </section>

                                {/* 2. Field Level Controls */}
                                <section className="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
                                    <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/30">
                                        <div>
                                            <h3 className="text-sm font-bold text-slate-900 flex items-center gap-2">
                                                Field Level Permission
                                            </h3>
                                            <p className="text-xs text-slate-500 mt-0.5">Fine-grained visibility and edit capability per column.</p>
                                        </div>
                                    </div>

                                    <div className="divide-y divide-slate-100">
                                        {resources.map(res => {
                                            const fields = fieldPermissions.filter(p => p.resource === res);
                                            if (!fields.length) return null;
                                            return (
                                                <div key={res} className="group">
                                                    <div className="px-6 py-3 bg-slate-50 border-b border-slate-100 flex justify-between items-center">
                                                        <span className="text-xs font-bold text-slate-800 uppercase tracking-widest">{res}</span>
                                                        <span className="text-[10px] font-mono text-slate-400">{fields.length} FIELDS</span>
                                                    </div>
                                                    <div className="px-6 py-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-x-8 gap-y-3">
                                                        {fields.map(field => (
                                                            <div key={field.field} className="flex items-center justify-between py-1">
                                                                <span className="text-sm text-slate-600 font-medium capitalize truncate pr-2" title={field.field}>
                                                                    {field.field.replace(/_/g, ' ')}
                                                                </span>
                                                                <div className="flex gap-1">
                                                                    <button
                                                                        onClick={() => handleFieldPermissionChange(res, field.field, 'view', !field.can_view)}
                                                                        className={`px-3 py-1 text-[10px] font-bold uppercase tracking-wider rounded transition-all ${field.can_view ? 'bg-slate-800 text-white' : 'bg-slate-100 text-slate-400 hover:bg-slate-200'
                                                                            }`}
                                                                    >
                                                                        View
                                                                    </button>
                                                                    <button
                                                                        onClick={() => handleFieldPermissionChange(res, field.field, 'edit', !field.can_edit)}
                                                                        className={`px-3 py-1 text-[10px] font-bold uppercase tracking-wider rounded transition-all ${field.can_edit ? 'bg-slate-800 text-white ring-1 ring-slate-800 ring-offset-1' : 'bg-slate-100 text-slate-400 hover:bg-slate-200'
                                                                            }`}
                                                                    >
                                                                        Edit
                                                                    </button>
                                                                </div>
                                                            </div>
                                                        ))}
                                                    </div>
                                                </div>
                                            )
                                        })}
                                    </div>
                                </section>

                                {/* 3. User Management (Simplified for this view) */}
                                {/* 3. User Management */}
                                <section className="pt-8 border-t border-slate-200">
                                    <div className="flex items-center justify-between mb-4">
                                        <h3 className="text-sm font-bold text-slate-900 uppercase tracking-wide">User Management</h3>
                                        <div className="flex gap-2">
                                            <input
                                                value={newUser.username}
                                                onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
                                                placeholder="New user email..."
                                                className="bg-white border border-slate-200 rounded px-2 py-1 text-xs focus:outline-none focus:border-slate-400 w-48"
                                            />
                                            <button onClick={async () => {
                                                if (!newUser.username) return toast.error("Email required");
                                                const tid = toast.loading("Inviting...");
                                                try {
                                                    await createUser({ email: newUser.username, role_id: selectedRole });
                                                    toast.success("Invited", { id: tid });
                                                    loadUsers();
                                                    setNewUser({ username: '', role_id: 1 });
                                                } catch (e) { toast.error("Failed", { id: tid }); }
                                            }} className="bg-slate-900 text-white px-3 py-1 rounded text-xs font-bold hover:bg-slate-800">Invite to Role</button>
                                        </div>
                                    </div>

                                    <div className="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
                                        <div className="px-6 py-3 bg-slate-50 border-b border-slate-100 grid grid-cols-12 text-xs font-semibold text-slate-500 uppercase tracking-wider">
                                            <div className="col-span-5">User</div>
                                            <div className="col-span-4">Role</div>
                                            <div className="col-span-2">Status</div>
                                            <div className="col-span-1 text-right">Actions</div>
                                        </div>
                                        <div className="divide-y divide-slate-100 max-h-60 overflow-y-auto">
                                            {users.map(u => (
                                                <div key={u.id} className="px-6 py-3 grid grid-cols-12 items-center hover:bg-slate-50 transition-colors group">
                                                    <div className="col-span-5 flex items-center gap-3">
                                                        <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-xs font-bold text-slate-500">
                                                            {u.username.charAt(0).toUpperCase()}
                                                        </div>
                                                        <div>
                                                            <div className="text-sm font-medium text-slate-900">{u.username}</div>
                                                            <div className="text-xs text-slate-400">{u.email}</div>
                                                        </div>
                                                    </div>
                                                    <div className="col-span-4">
                                                        {u.username === 'admin' ? (
                                                            <span className="text-xs font-medium text-slate-500 bg-slate-100 px-2 py-1 rounded">Admin System</span>
                                                        ) : (
                                                            <select
                                                                value={u.role_id}
                                                                onChange={async (e) => {
                                                                    try {
                                                                        await updateUserRole(u.id, Number(e.target.value));
                                                                        toast.success(`Updated ${u.username}'s role`);
                                                                        loadUsers();
                                                                    } catch (err) {
                                                                        toast.error("Failed to update role");
                                                                    }
                                                                }}
                                                                className="w-full text-xs border-none bg-transparent focus:ring-0 text-slate-700 font-medium cursor-pointer hover:bg-slate-100 rounded px-1 -ml-1 py-1"
                                                            >
                                                                {roles.map(r => (
                                                                    <option key={r.id} value={r.id}>{r.name}</option>
                                                                ))}
                                                            </select>
                                                        )}
                                                    </div>
                                                    <div className="col-span-2">
                                                        <span className={`text-[10px] px-2 py-0.5 rounded-full font-medium border ${u.status === 'Active' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-amber-50 text-amber-600 border-amber-100'}`}>
                                                            {u.status || 'Active'}
                                                        </span>
                                                    </div>
                                                    <div className="col-span-1 text-right opacity-0 group-hover:opacity-100 transition-opacity">
                                                        {u.username !== 'admin' && (
                                                            <button
                                                                onClick={() => setDeleteConfirmation({ isOpen: true, userId: u.id, username: u.username, type: 'user' })}
                                                                className="text-slate-400 hover:text-red-600 p-1"
                                                                title="Remove User"
                                                            >
                                                                <Trash2 className="w-4 h-4" />
                                                            </button>
                                                        )}
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </section>
                            </div>
                        </div>
                    ) : (
                        <div className="flex-1 flex flex-col items-center justify-center text-slate-400">
                            <div className="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mb-4">
                                <span className="text-2xl">👋</span>
                            </div>
                            <h3 className="text-lg font-medium text-slate-900">Select a Role</h3>
                            <p className="max-w-xs text-center text-sm mt-2">Choose a role from the sidebar to configure granular access permissions.</p>
                        </div>
                    )}
                </main>
            </div>

            {/* Delete Confirmation Modal reused similarly */}
            <AnimatePresence>
                {deleteConfirmation.isOpen && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/20 backdrop-blur-sm">
                        <motion.div initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }} className="bg-white rounded-lg shadow-xl p-6 max-w-sm w-full">
                            <h3 className="text-lg font-bold text-slate-900 mb-2">Confirm Deletion</h3>
                            <p className="text-sm text-slate-500 mb-6">Permanently remove <strong className="text-slate-900">{deleteConfirmation.username}</strong>?</p>
                            <div className="flex gap-3">
                                <button onClick={() => setDeleteConfirmation({ isOpen: false, userId: null, username: '', type: 'user' })} className="flex-1 py-2 bg-slate-100 text-slate-700 rounded font-medium hover:bg-slate-200">Cancel</button>
                                <button onClick={handleDelete} className="flex-1 py-2 bg-red-600 text-white rounded font-medium hover:bg-red-700">Delete</button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </div>
    );
}
