
import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search } from "lucide-react";
import CreateUserForm from "@/components/users/CreateUserForm";

type User = {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  role: string;
  status: string;
  lastActive: string;
};

const mockUsers: User[] = [
  { id: 1, firstName: "John", lastName: "Doe", email: "john@example.com", role: "Admin", status: "active", lastActive: "2024-02-20" },
  { id: 2, firstName: "Jane", lastName: "Smith", email: "jane@example.com", role: "Editor", status: "active", lastActive: "2024-02-19" },
  { id: 3, firstName: "Mike", lastName: "Johnson", email: "mike@example.com", role: "Viewer", status: "inactive", lastActive: "2024-02-15" },
  { id: 4, firstName: "Sarah", lastName: "Williams", email: "sarah@example.com", role: "Editor", status: "active", lastActive: "2024-02-14" },
  { id: 5, firstName: "David", lastName: "Brown", email: "david@example.com", role: "Viewer", status: "active", lastActive: "2024-02-13" },
  { id: 6, firstName: "Lisa", lastName: "Anderson", email: "lisa@example.com", role: "Admin", status: "active", lastActive: "2024-02-12" },
  { id: 7, firstName: "Robert", lastName: "Taylor", email: "robert@example.com", role: "Editor", status: "inactive", lastActive: "2024-02-11" },
  { id: 8, firstName: "Emily", lastName: "Jones", email: "emily@example.com", role: "Viewer", status: "active", lastActive: "2024-02-10" },
  { id: 9, firstName: "Michael", lastName: "Wilson", email: "michael@example.com", role: "Editor", status: "active", lastActive: "2024-02-09" },
  { id: 10, firstName: "Jennifer", lastName: "Thomas", email: "jennifer@example.com", role: "Viewer", status: "active", lastActive: "2024-02-08" },
  { id: 11, firstName: "Daniel", lastName: "Martinez", email: "daniel@example.com", role: "Admin", status: "active", lastActive: "2024-02-07" },
  { id: 12, firstName: "Maria", lastName: "Garcia", email: "maria@example.com", role: "Editor", status: "active", lastActive: "2024-02-06" },
  { id: 13, firstName: "William", lastName: "White", email: "william@example.com", role: "Viewer", status: "inactive", lastActive: "2024-02-05" },
  { id: 14, firstName: "Elizabeth", lastName: "Lee", email: "elizabeth@example.com", role: "Editor", status: "active", lastActive: "2024-02-04" },
  { id: 15, firstName: "Richard", lastName: "Clark", email: "richard@example.com", role: "Viewer", status: "active", lastActive: "2024-02-03" }
];

const Users = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  const filteredUsers = mockUsers.filter(
    (user) =>
      user.firstName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.lastName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.email.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const totalPages = Math.ceil(filteredUsers.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const displayedUsers = filteredUsers.slice(
    startIndex,
    startIndex + itemsPerPage
  );

  return (
    <div className="container animate-fadeIn">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-semibold">Users</h1>
        <CreateUserForm />
      </div>

      <Card className="p-6 glass mb-6">
        <div className="flex gap-4 mb-6">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <Input
              placeholder="Search users..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9"
            />
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="text-left border-b">
                <th className="pb-3">Name</th>
                <th className="pb-3">Email</th>
                <th className="pb-3">Role</th>
                <th className="pb-3">Status</th>
                <th className="pb-3">Last Active</th>
              </tr>
            </thead>
            <tbody>
              {displayedUsers.map((user) => (
                <tr key={user.id} className="border-b last:border-0">
                  <td className="py-3">
                    {user.firstName} {user.lastName}
                  </td>
                  <td className="py-3">{user.email}</td>
                  <td className="py-3">{user.role}</td>
                  <td className="py-3">
                    <span
                      className={`px-2 py-1 rounded-full text-xs ${
                        user.status === "active"
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {user.status}
                    </span>
                  </td>
                  <td className="py-3">{user.lastActive}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="flex justify-between items-center mt-4">
          <p className="text-sm text-muted-foreground">
            Showing {startIndex + 1} to{" "}
            {Math.min(startIndex + itemsPerPage, filteredUsers.length)} of{" "}
            {filteredUsers.length} entries
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => setCurrentPage(currentPage - 1)}
              disabled={currentPage === 1}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              onClick={() => setCurrentPage(currentPage + 1)}
              disabled={currentPage === totalPages}
            >
              Next
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default Users;
