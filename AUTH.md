# 🔐 Authentication & Authorization Guide

## 🎯 **Overview**

The dental clinic system implements **JWT-based authentication** with **Role-Based Access Control (RBAC)** to ensure secure access to different system features.

## 👥 **User Roles & Permissions**

### **🔴 Admin**
- **Full system access**
- Can create/delete any user
- Can perform all CRUD operations
- System administration privileges

### **🟢 Doctor**
- Can **complete appointments** (add medical data)
- Can **CRUD formulas** (dental records)
- Can **CRUD statuses** (diagnosis, treatment, tooth conditions)
- Can **CRUD complaints**
- Can **view all user details** (for patient cards)
- Cannot create/delete appointments

### **🟡 Receptionist**
- Can **create/update/delete appointments**
- Can **view client information**
- Can **view doctors** for appointment booking
- Cannot access medical data (formulas)
- Cannot complete appointments

### **🔵 Client**
- **Cannot authenticate** (no login access)
- Created by staff members only
- No password/email required

## 🚀 **Authentication Endpoints**

### **1. Sign Up (Staff Registration)**
```http
POST /api/v1/auth/signup
Content-Type: application/json

{
  "email": "dr.smith@clinic.com",
  "password": "securepassword123",
  "full_name": "Dr. John Smith",
  "role": "doctor",
  "phone_number": "+1234567890"
}
```

**Response:**
```json
{
  "message": "User created successfully",
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

### **2. Sign In**
```http
POST /api/v1/auth/signin
Content-Type: application/json

{
  "email": "dr.smith@clinic.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "message": "Sign in successful",
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

### **3. Refresh Token**
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### **4. Get Profile**
```http
GET /api/v1/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response:**
```json
{
  "user": {
    "user_id": "507f1f77bcf86cd799439011",
    "email": "dr.smith@clinic.com",
    "full_name": "Dr. John Smith",
    "role": "doctor"
  }
}
```

## 🔒 **Protected Routes & Permissions**

### **User Management**
| Method | Endpoint | Admin | Doctor | Receptionist |
|--------|----------|-------|--------|--------------|
| POST | `/users` | ✅ | ❌ | ❌ |
| GET | `/users` | ✅ | ❌ | ❌ |
| GET | `/users/:id` | ✅ | ✅ | ✅ |
| PUT | `/users/:id` | ✅ | ✅ | ✅ |
| DELETE | `/users/:id` | ✅ | ❌ | ❌ |
| GET | `/users/doctors` | ✅ | ✅ | ✅ |
| GET | `/users/clients` | ✅ | ✅ | ✅ |

### **Appointment Management**
| Method | Endpoint | Admin | Doctor | Receptionist |
|--------|----------|-------|--------|--------------|
| POST | `/appointments` | ✅ | ✅ | ✅ |
| GET | `/appointments` | ✅ | ✅ | ✅ |
| GET | `/appointments/:id` | ✅ | ✅ | ✅ |
| PUT | `/appointments/:id` | ✅ | ✅ | ✅ |
| DELETE | `/appointments/:id` | ✅ | ✅ | ✅ |
| POST | `/appointments/:id/complete` | ✅ | ✅ | ❌ |
| POST | `/appointments/:id/cancel` | ✅ | ✅ | ✅ |

### **Medical Data (Formulas)**
| Method | Endpoint | Admin | Doctor | Receptionist |
|--------|----------|-------|--------|--------------|
| GET | `/formulas/:id` | ✅ | ✅ | ❌ |
| GET | `/formulas/user/:userId` | ✅ | ✅ | ❌ |

### **Statuses & Complaints**
| Method | Endpoint | Admin | Doctor | Receptionist |
|--------|----------|-------|--------|--------------|
| POST | `/statuses` | ✅ | ✅ | ❌ |
| GET | `/statuses/*` | ✅ | ✅ | ✅ |
| PUT | `/statuses/:id` | ✅ | ✅ | ❌ |
| DELETE | `/statuses/:id` | ✅ | ✅ | ❌ |
| POST | `/complaints` | ✅ | ✅ | ❌ |
| GET | `/complaints/*` | ✅ | ✅ | ✅ |

## 🛠️ **Frontend Integration**

### **1. Authentication Flow**

```javascript
// Sign up a new doctor
const signUp = async () => {
  const response = await fetch('/api/v1/auth/signup', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: 'dr.smith@clinic.com',
      password: 'securepassword123',
      full_name: 'Dr. John Smith',
      role: 'doctor'
    })
  });
  
  const data = await response.json();
  
  // Store tokens
  localStorage.setItem('access_token', data.tokens.access_token);
  localStorage.setItem('refresh_token', data.tokens.refresh_token);
};

// Sign in existing user
const signIn = async () => {
  const response = await fetch('/api/v1/auth/signin', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: 'dr.smith@clinic.com',
      password: 'securepassword123'
    })
  });
  
  const data = await response.json();
  localStorage.setItem('access_token', data.tokens.access_token);
  localStorage.setItem('refresh_token', data.tokens.refresh_token);
};
```

### **2. Making Authenticated Requests**

```javascript
// Function to get access token
const getAccessToken = () => localStorage.getItem('access_token');

// Function to make authenticated requests
const authenticatedFetch = async (url, options = {}) => {
  const token = getAccessToken();
  
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
    ...options.headers
  };
  
  const response = await fetch(url, {
    ...options,
    headers
  });
  
  // Handle token expiration
  if (response.status === 401) {
    await refreshToken();
    // Retry request with new token
    return authenticatedFetch(url, options);
  }
  
  return response;
};

// Example: Create appointment (receptionist/doctor)
const createAppointment = async (appointmentData) => {
  return authenticatedFetch('/api/v1/appointments', {
    method: 'POST',
    body: JSON.stringify(appointmentData)
  });
};

// Example: Complete appointment (doctor only)
const completeAppointment = async (appointmentId, medicalData) => {
  return authenticatedFetch(`/api/v1/appointments/${appointmentId}/complete`, {
    method: 'POST',
    body: JSON.stringify(medicalData)
  });
};
```

### **3. Token Refresh Implementation**

```javascript
const refreshToken = async () => {
  const refreshToken = localStorage.getItem('refresh_token');
  
  try {
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });
    
    if (response.ok) {
      const data = await response.json();
      localStorage.setItem('access_token', data.tokens.access_token);
      localStorage.setItem('refresh_token', data.tokens.refresh_token);
      return true;
    } else {
      // Refresh failed, redirect to login
      logout();
      return false;
    }
  } catch (error) {
    logout();
    return false;
  }
};

const logout = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  window.location.href = '/login';
};
```

### **4. Role-Based UI Components**

```jsx
// React component with role-based rendering
const AppointmentCard = ({ appointment, userRole }) => {
  return (
    <div className="appointment-card">
      <h3>{appointment.client_name}</h3>
      <p>{appointment.date_time}</p>
      
      {/* Only doctors can complete appointments */}
      {userRole === 'doctor' && appointment.status === 'scheduled' && (
        <button onClick={() => completeAppointment(appointment.id)}>
          Complete Appointment
        </button>
      )}
      
      {/* Receptionists and doctors can modify appointments */}
      {(userRole === 'receptionist' || userRole === 'doctor') && (
        <button onClick={() => editAppointment(appointment.id)}>
          Edit Appointment
        </button>
      )}
      
      {/* Only doctors can view dental records */}
      {userRole === 'doctor' && (
        <button onClick={() => viewDentalRecord(appointment.client_id)}>
          View Dental Record
        </button>
      )}
    </div>
  );
};
```

## 🔐 **Security Features**

### **JWT Token Structure**
```json
{
  "user_id": "507f1f77bcf86cd799439011",
  "email": "dr.smith@clinic.com",
  "full_name": "Dr. John Smith",
  "role": "doctor",
  "type": "access",
  "exp": 1640995200,
  "iat": 1640908800
}
```

### **Password Security**
- **Bcrypt hashing** with cost factor 12
- **Minimum 6 characters** required
- **Salt included** in hash

### **Token Management**
- **Access tokens**: Short-lived (configurable, default 24h)
- **Refresh tokens**: Long-lived (7x access token duration)
- **Automatic refresh** on expiration
- **Secure storage** recommended (httpOnly cookies in production)

## 🚨 **Error Handling**

### **Common Error Responses**

```json
// Invalid credentials
{
  "error": "invalid email or password"
}

// Insufficient permissions
{
  "error": "Insufficient permissions"
}

// Token expired
{
  "error": "Invalid or expired token"
}

// Role validation failed
{
  "error": "invalid role"
}
```

## 🎯 **Example Workflow**

### **Daily Clinic Operations**

1. **Morning Setup (Receptionist)**
   ```bash
   # Sign in
   POST /auth/signin
   
   # View today's appointments
   GET /appointments?date=2024-01-15
   
   # Create new appointment
   POST /appointments
   ```

2. **Patient Visit (Doctor)**
   ```bash
   # View appointment details
   GET /appointments/507f1f77bcf86cd799439011
   
   # Complete appointment with medical data
   POST /appointments/507f1f77bcf86cd799439011/complete
   
   # Update dental formula
   GET /formulas/user/507f1f77bcf86cd799439012
   ```

3. **Administrative Tasks (Admin)**
   ```bash
   # Add new doctor
   POST /auth/signup
   
   # View all users
   GET /users
   
   # System maintenance
   GET /statuses
   POST /statuses
   ```

Your authentication system is now fully implemented with comprehensive RBAC! 🦷🔐