db.createUser({
  user: "admin",
  pwd: "password123",
  roles: [{ role: "readWrite", db: "hospital" }],
});

db = new Mongo().getDB("hospital");

db.patients.insertMany([
  { firstName: "John", lastName: "Doe", dob: "1990-01-01", gender: "Male" },
  { firstName: "Jane", lastName: "Doe", dob: "1995-02-01", gender: "Female" },
]);

db.users.insertOne({
  username: "admin",
  password: "password", // Use a hashed password for production
});
