package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Patient represents a patient document in MongoDB
type Patient struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	DOB       time.Time `bson:"dob"`
	Phone     string    `bson:"phone"`
	State     string    `bson:"state"`
	CreatedAt time.Time `bson:"created_at"`
}

// Address represents an address document in MongoDB
type Address struct {
	ID        string `bson:"_id"`
	PatientID string `bson:"patient_id"`
	Line1     string `bson:"line1"`
	Line2     string `bson:"line2"`
	City      string `bson:"city"`
	State     string `bson:"state"`
	Zip       string `bson:"zip"`
}

// Prescription represents a prescription document in MongoDB
type Prescription struct {
	ID        string    `bson:"_id"`
	PatientID string    `bson:"patient_id"`
	Drug      string    `bson:"drug"`
	Dose      string    `bson:"dose"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
}

func main() {
	fmt.Println("🌱 Starting MongoDB seeding...")

	// MongoDB connection string
	uri := "mongodb://admin:admin123@localhost:27017"
	dbName := "pharmacy_modernization"

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal("Failed to disconnect:", err)
		}
	}()

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("✅ Connected to MongoDB")

	db := client.Database(dbName)

	// Seed patients
	seedPatients(db)

	// Seed addresses
	seedAddresses(db)

	// Seed prescriptions
	seedPrescriptions(db)

	fmt.Println("\n🎉 Seeding complete! You can view the data at:")
	fmt.Println("   MongoDB: mongodb://admin:admin123@localhost:27017/pharmacy_modernization")
	fmt.Println("   Mongo Express: http://localhost:8081")
}

func seedPatients(db *mongo.Database) {
	collection := db.Collection("patients")

	fmt.Println("\n👥 Seeding patients...")
	fmt.Println("🗑️  Clearing existing patients...")
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Failed to clear patients collection:", err)
	}

	patients := []Patient{
		{
			ID:        "P001",
			Name:      "Ava Thompson2",
			DOB:       time.Date(1988, time.January, 12, 0, 0, 0, 0, time.UTC),
			Phone:     "(206) 417-8842",
			State:     "Washington",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P002",
			Name:      "Liam Anderson",
			DOB:       time.Date(1979, time.March, 3, 0, 0, 0, 0, time.UTC),
			Phone:     "(415) 736-5528",
			State:     "California",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P003",
			Name:      "Sophia Martinez",
			DOB:       time.Date(1992, time.July, 27, 0, 0, 0, 0, time.UTC),
			Phone:     "(617) 980-3314",
			State:     "Massachusetts",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P004",
			Name:      "Noah Patel",
			DOB:       time.Date(1985, time.May, 5, 0, 0, 0, 0, time.UTC),
			Phone:     "(972) 645-2091",
			State:     "Texas",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P005",
			Name:      "Mia Chen",
			DOB:       time.Date(1996, time.September, 19, 0, 0, 0, 0, time.UTC),
			Phone:     "(312) 478-6605",
			State:     "Illinois",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P006",
			Name:      "Ethan Johnson",
			DOB:       time.Date(1975, time.November, 8, 0, 0, 0, 0, time.UTC),
			Phone:     "(303) 825-1947",
			State:     "Colorado",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P007",
			Name:      "Olivia Rossi",
			DOB:       time.Date(1990, time.February, 22, 0, 0, 0, 0, time.UTC),
			Phone:     "(646) 291-0743",
			State:     "New York",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P008",
			Name:      "Jackson Lee",
			DOB:       time.Date(1983, time.April, 16, 0, 0, 0, 0, time.UTC),
			Phone:     "(503) 913-2286",
			State:     "Oregon",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P009",
			Name:      "Emma Davis",
			DOB:       time.Date(1998, time.December, 2, 0, 0, 0, 0, time.UTC),
			Phone:     "(305) 744-1189",
			State:     "Florida",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P010",
			Name:      "Lucas Hernandez",
			DOB:       time.Date(1981, time.June, 14, 0, 0, 0, 0, time.UTC),
			Phone:     "(713) 402-5378",
			State:     "Texas",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P011",
			Name:      "Isabella Rodriguez",
			DOB:       time.Date(1994, time.August, 30, 0, 0, 0, 0, time.UTC),
			Phone:     "(212) 555-9876",
			State:     "New York",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P012",
			Name:      "Mason Williams",
			DOB:       time.Date(1987, time.October, 18, 0, 0, 0, 0, time.UTC),
			Phone:     "(404) 555-4321",
			State:     "Georgia",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P013",
			Name:      "Charlotte Brown",
			DOB:       time.Date(1991, time.April, 7, 0, 0, 0, 0, time.UTC),
			Phone:     "(602) 555-1111",
			State:     "Arizona",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P014",
			Name:      "James Taylor",
			DOB:       time.Date(1980, time.December, 25, 0, 0, 0, 0, time.UTC),
			Phone:     "(702) 555-2222",
			State:     "Nevada",
			CreatedAt: time.Now(),
		},
		{
			ID:        "P015",
			Name:      "Amelia Garcia",
			DOB:       time.Date(1995, time.June, 15, 0, 0, 0, 0, time.UTC),
			Phone:     "(214) 555-3333",
			State:     "Texas",
			CreatedAt: time.Now(),
		},
	}

	docs := make([]interface{}, len(patients))
	for i, p := range patients {
		docs[i] = p
	}

	result, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		log.Fatal("Failed to insert patients:", err)
	}

	fmt.Printf("✅ Successfully inserted %d patients!\n", len(result.InsertedIDs))
}

func seedAddresses(db *mongo.Database) {
	collection := db.Collection("addresses")

	fmt.Println("\n🏠 Seeding addresses...")
	fmt.Println("🗑️  Clearing existing addresses...")
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Failed to clear addresses collection:", err)
	}

	addresses := []Address{
		// Ava Thompson (P001)
		{ID: "A001", PatientID: "P001", Line1: "123 Main St", Line2: "Apt 4B", City: "Seattle", State: "WA", Zip: "98101"},
		{ID: "A002", PatientID: "P001", Line1: "456 Market Ave", Line2: "", City: "Seattle", State: "WA", Zip: "98102"},

		// Liam Anderson (P002)
		{ID: "A003", PatientID: "P002", Line1: "789 Sunset Blvd", Line2: "", City: "San Francisco", State: "CA", Zip: "94102"},

		// Sophia Martinez (P003)
		{ID: "A004", PatientID: "P003", Line1: "321 Commonwealth Ave", Line2: "Unit 12", City: "Boston", State: "MA", Zip: "02215"},

		// Noah Patel (P004)
		{ID: "A005", PatientID: "P004", Line1: "555 Ranch Road", Line2: "", City: "Dallas", State: "TX", Zip: "75201"},
		{ID: "A006", PatientID: "P004", Line1: "888 Business Park Dr", Line2: "Suite 300", City: "Dallas", State: "TX", Zip: "75202"},

		// Mia Chen (P005)
		{ID: "A007", PatientID: "P005", Line1: "999 Lake Shore Dr", Line2: "", City: "Chicago", State: "IL", Zip: "60611"},

		// Ethan Johnson (P006)
		{ID: "A008", PatientID: "P006", Line1: "777 Mountain View Rd", Line2: "", City: "Denver", State: "CO", Zip: "80202"},

		// Olivia Rossi (P007)
		{ID: "A009", PatientID: "P007", Line1: "222 Broadway", Line2: "Floor 15", City: "New York", State: "NY", Zip: "10007"},
		{ID: "A010", PatientID: "P007", Line1: "333 Park Ave", Line2: "", City: "New York", State: "NY", Zip: "10022"},

		// Jackson Lee (P008)
		{ID: "A011", PatientID: "P008", Line1: "444 Forest Lane", Line2: "", City: "Portland", State: "OR", Zip: "97201"},

		// Emma Davis (P009)
		{ID: "A012", PatientID: "P009", Line1: "666 Ocean Drive", Line2: "Apt 23", City: "Miami", State: "FL", Zip: "33139"},

		// Lucas Hernandez (P010)
		{ID: "A013", PatientID: "P010", Line1: "111 Heritage St", Line2: "", City: "Houston", State: "TX", Zip: "77002"},
	}

	docs := make([]interface{}, len(addresses))
	for i, a := range addresses {
		docs[i] = a
	}

	result, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		log.Fatal("Failed to insert addresses:", err)
	}

	fmt.Printf("✅ Successfully inserted %d addresses!\n", len(result.InsertedIDs))
}

func seedPrescriptions(db *mongo.Database) {
	collection := db.Collection("prescriptions")

	fmt.Println("\n💊 Seeding prescriptions...")
	fmt.Println("🗑️  Clearing existing prescriptions...")
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Failed to clear prescriptions collection:", err)
	}

	now := time.Now()
	prescriptions := []Prescription{
		// Ava Thompson (P001)
		{ID: "RX001", PatientID: "P001", Drug: "Lisinopril", Dose: "10mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -30)},
		{ID: "RX002", PatientID: "P001", Drug: "Metformin", Dose: "500mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -25)},

		// Liam Anderson (P002)
		{ID: "RX003", PatientID: "P002", Drug: "Atorvastatin", Dose: "20mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -45)},
		{ID: "RX004", PatientID: "P002", Drug: "Omeprazole", Dose: "40mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -20)},
		{ID: "RX005", PatientID: "P002", Drug: "Amoxicillin", Dose: "500mg", Status: "Completed", CreatedAt: now.AddDate(0, 0, -60)},

		// Sophia Martinez (P003)
		{ID: "RX006", PatientID: "P003", Drug: "Albuterol", Dose: "90mcg", Status: "Active", CreatedAt: now.AddDate(0, 0, -15)},
		{ID: "RX007", PatientID: "P003", Drug: "Montelukast", Dose: "10mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -10)},

		// Noah Patel (P004)
		{ID: "RX008", PatientID: "P004", Drug: "Losartan", Dose: "50mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -40)},
		{ID: "RX009", PatientID: "P004", Drug: "Amlodipine", Dose: "5mg", Status: "Paused", CreatedAt: now.AddDate(0, 0, -35)},

		// Mia Chen (P005)
		{ID: "RX010", PatientID: "P005", Drug: "Sertraline", Dose: "50mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -90)},
		{ID: "RX011", PatientID: "P005", Drug: "Ibuprofen", Dose: "400mg", Status: "Completed", CreatedAt: now.AddDate(0, 0, -100)},

		// Ethan Johnson (P006)
		{ID: "RX012", PatientID: "P006", Drug: "Warfarin", Dose: "5mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -120)},
		{ID: "RX013", PatientID: "P006", Drug: "Furosemide", Dose: "40mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -115)},
		{ID: "RX014", PatientID: "P006", Drug: "Levothyroxine", Dose: "75mcg", Status: "Active", CreatedAt: now.AddDate(0, 0, -110)},

		// Olivia Rossi (P007)
		{ID: "RX015", PatientID: "P007", Drug: "Escitalopram", Dose: "10mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -55)},

		// Jackson Lee (P008)
		{ID: "RX016", PatientID: "P008", Drug: "Gabapentin", Dose: "300mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -70)},
		{ID: "RX017", PatientID: "P008", Drug: "Tramadol", Dose: "50mg", Status: "Paused", CreatedAt: now.AddDate(0, 0, -65)},

		// Emma Davis (P009)
		{ID: "RX018", PatientID: "P009", Drug: "Azithromycin", Dose: "250mg", Status: "Completed", CreatedAt: now.AddDate(0, 0, -5)},
		{ID: "RX019", PatientID: "P009", Drug: "Birth Control", Dose: "Daily", Status: "Active", CreatedAt: now.AddDate(0, 0, -180)},

		// Lucas Hernandez (P010)
		{ID: "RX020", PatientID: "P010", Drug: "Simvastatin", Dose: "40mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -200)},
		{ID: "RX021", PatientID: "P010", Drug: "Aspirin", Dose: "81mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -195)},
		{ID: "RX022", PatientID: "P010", Drug: "Clopidogrel", Dose: "75mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -190)},

		// Additional prescriptions for variety
		{ID: "RX023", PatientID: "P011", Drug: "Prednisone", Dose: "20mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -12)},
		{ID: "RX024", PatientID: "P012", Drug: "Metoprolol", Dose: "50mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -80)},
		{ID: "RX025", PatientID: "P013", Drug: "Hydrochlorothiazide", Dose: "25mg", Status: "Active", CreatedAt: now.AddDate(0, 0, -50)},
		{ID: "RX026", PatientID: "P014", Drug: "Insulin Glargine", Dose: "20 units", Status: "Active", CreatedAt: now.AddDate(0, 0, -150)},
		{ID: "RX027", PatientID: "P015", Drug: "Rosuvastatin", Dose: "10mg", Status: "Draft", CreatedAt: now.AddDate(0, 0, -2)},
	}

	docs := make([]interface{}, len(prescriptions))
	for i, p := range prescriptions {
		docs[i] = p
	}

	result, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		log.Fatal("Failed to insert prescriptions:", err)
	}

	fmt.Printf("✅ Successfully inserted %d prescriptions!\n", len(result.InsertedIDs))
}
