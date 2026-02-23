import React, { useState, useEffect, ChangeEvent } from 'react';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { TextField, Button, Box } from '@mui/material';

import { fetchUsers, User } from './services/fetchUsers';

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', width: 100 },
  { field: 'email', headerName: 'Email', width: 200 },
  { field: 'username', headerName: 'Username', width: 200 },
  { field: 'domain', headerName: 'Domain', width: 200 },
  { field: 'password', headerName: 'Password', width: 200 },
];

function App() {
  const [users, setUsers] = useState<User[]>([]);
  const [emailFilter, setEmailFilter] = useState('');
  const [domainFilter, setDomainFilter] = useState('');
  const [loading, setLoading] = useState(false);

  const getUsers = async (email?: string, domain?: string) => {
    try {
      setLoading(true);
      const fetchedUsers = await fetchUsers(email, domain);
      setUsers(fetchedUsers);
    } catch (error) {
      console.error('Error fetching users:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getUsers();
  }, []);

  const handleFilter = () => {
    getUsers(emailFilter, domainFilter);
  };

  const handleEmailChange = (event: ChangeEvent<HTMLInputElement>) => {
    setEmailFilter(event.target.value);
  };

  const handleDomainChange = (event: ChangeEvent<HTMLInputElement>) => {
    setDomainFilter(event.target.value);
  };

  return (
    <Box sx={{ height: 600, width: '90%', p: 2 }}>
      <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
        <TextField
          label="Filter by Email"
          variant="outlined"
          value={emailFilter}
          onChange={handleEmailChange}
        />
        <TextField
          label="Filter by Domain"
          variant="outlined"
          value={domainFilter}
          onChange={handleDomainChange}
        />
        <Button variant="contained" onClick={handleFilter} disabled={loading}>
          Apply Filters
        </Button>
      </Box>

      <DataGrid
        rows={users}
        columns={columns}
        loading={loading}
        getRowId={(row) => row.id}
      />
    </Box>
  );
}

export default App;
