import axios from 'axios';

const API_URL = 'http://localhost:5001/api';

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const login = async (username, password) => {
    const response = await api.post('/auth/login', { username, password });
    return response.data;
};

export const declineInvitation = async (token) => {
    const response = await api.post('/auth/decline-invitation', { token });
    return response.data;
};

export const fetchMyPermissions = async () => {
    const response = await api.get('/auth/permissions');
    return response.data;
};

// Data Resource APIs
export const fetchResource = async (resource) => {
    const response = await api.get(`/data/${resource}`);
    return response.data;
};

export const createResource = async (resource, data) => {
    const response = await api.post(`/data/${resource}`, data);
    return response.data;
};

export const updateResource = async (resource, id, data) => {
    const response = await api.put(`/data/${resource}/${id}`, data);
    return response.data;
};

export const deleteResource = async (resource, id) => {
    const response = await api.delete(`/data/${resource}/${id}`);
    return response.data;
};

// Admin APIs
export const fetchRoles = async () => {
    const response = await api.get('/admin/roles');
    return response.data;
};

export const createRole = async (name) => {
    const response = await api.post('/admin/roles', { name });
    return response.data;
};

export const fetchPermissions = async (roleId) => {
    const response = await api.get(`/admin/permissions/${roleId}`);
    return response.data;
};

export const savePermissions = async (permissionData) => {
    try {
        const response = await api.post('/admin/permissions', permissionData);
        return response.data;
    } catch (err) {
        console.error("Save failed", err);
        throw err;
    }
};

export const deletePermission = async (perm) => {
    const response = await api.delete('/admin/permissions', { params: perm });
    return response.data;
};

export const fetchUsers = async () => {
    const response = await api.get('/admin/users');
    return response.data;
};


export const createUser = async (userData) => {
    const response = await api.post('/admin/users', userData);
    return response.data;
};

export const updateUserRole = async (userId, roleId) => {
    const response = await api.put(`/admin/users/${userId}/role`, { role_id: roleId });
    return response.data;
};

export const deleteUser = async (userId) => {
    const response = await api.delete(`/admin/users/${userId}`);
    return response.data;
};

// Available to all authenticated users
export const fetchAllUsers = async () => {
    const response = await api.get('/auth/users');
    return response.data;
};

// Fetch employee names for dropdowns (authenticated)
export const fetchEmployeeNames = async () => {
    const response = await api.get('/auth/employees');
    return response.data;
};
