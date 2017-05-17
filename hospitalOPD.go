package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"strconv"
)

type SimpleChaincode struct {
}

type POLICY_ID_Holder struct {
	POLICY_IDs []string `json:"POLICY_IDs"`
}
type HOSPITAL_NAME_Holder struct {
	HOSPITAL_NAMEs []string `json:"HOSPITAL_NAMEs"`
}
type PATIENT_HISTORY struct {
	PROPERTY_HISTORY_IDs []string `json:"PROPERTY_HISTORY_IDs"`
}

//Information to be stored about land in blockchain network
type Patient struct {
	PloicyId         string `json:"ploicyId"`
	City             string `json:"city"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Contact_Number   int64  `json:"contact_Number"`
	Hospital         string `json:"hospital"`
	AppointmentTime  string `json:"appointmentTime"`
	UnavailedBalance int64  `json:"unavailedBalance"`
	ClaimedAmount    int64  `json:"claimedAmount"`
}

type Hospital struct {
	PloicyId         string `json:"ploicyId"`
	Hospital         string `json:"hospital"`
	AppointmentTime  string `json:"appointmentTime"`
	UnavailedBalance int64  `json:"unavailedBalance"`
	ClaimedAmount    int64  `json:"claimedAmount"`
}

type PatientyByHospital struct {
	hospitalDetails Hospital  `json:"ownerDetails"` //generated by blockchain
	PatientDetails  []Patient `json:"propertyDetails"`
}

/*type PropertyHistory struct {
	HistoryId       int    `json:"historyId"`
	PropertyId      string `json:"propertyId"` //generated by blockchain
	OwnerId         string `json:"ownerId"`
	AgreementDate   string `json:"agreementDate"`
	AgreementAmount string `json:"agreementAmount"`
}*/

func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)

	}
}

//Init
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	var policyIdBytes []byte
	var hospitalNameBytes []byte

	var policyIds POLICY_ID_Holder
	var hospitalNames HOSPITAL_NAME_Holder

	policyIdBytes, err = json.Marshal(policyIds)
	hospitalNameBytes, err = json.Marshal(hospitalNames)

	if err != nil {
		return nil, errors.New("Error creating OWNER_ID_Holder record")
	}

	err = stub.PutState("policy_Ids", policyIdBytes)
	err = stub.PutState("hospital_Names", hospitalNameBytes)

	return nil, nil

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "createAppointment" {

		return createAppointment(stub, args)
	}
	if function == "updateBalanceAPI" {

		return updateBalanceAPI(stub, args)
	}

	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "listAppointmentsByHospital" {

		return listAppointmentsByHospital(stub, args)
	}
	if function == "listAllAppointments" {

		return listAllAppointments(stub, args)
	}

	return nil, nil

}

func listAppointmentsByHospital(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	hospitalName := args[0]

	bytes, err := stub.GetState("policy_Ids")

	fmt.Println("Ids recieved", string(bytes))
	var policyIdHolder POLICY_ID_Holder
	err = json.Unmarshal(bytes, &policyIdHolder)

	result := []Patient{}

	var p Patient
	var resultBytes []byte

	for _, pat := range policyIdHolder.POLICY_IDs {

		fmt.Println("Inside for loop for getting PatientDetails for Hospital. Policy Id is  ", pat)

		p, err = retrievePatient(stub, pat)

		if p.Hospital == hospitalName {
			//temp, err = json.Marshal(p)

			if err == nil {

				result = append(result, p)
			}

		}

	}

	resultBytes, err = json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return resultBytes, nil

}

func listAllAppointments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inside list appointed Patients")
	var resultBytes []byte

	bytes, err := stub.GetState("policy_Ids")

	fmt.Println("Ids recieved", string(bytes))
	var policyIdHolder POLICY_ID_Holder
	err = json.Unmarshal(bytes, &policyIdHolder)

	result := []Patient{}

	var p Patient

	for _, pat := range policyIdHolder.POLICY_IDs {

		fmt.Println("Inside for loop for getting Property. Policy  Id is  ", pat)

		p, err = retrievePatient(stub, pat)

		if err == nil {
			result = append(result, p)
		}

	}
	resultBytes, err = json.Marshal(result)

	return resultBytes, nil
}

func retrievePatient(stub shim.ChaincodeStubInterface, patientId string) (Patient, error) {

	fmt.Println("Inside retrieve Patient")

	var p Patient

	bytes, err := stub.GetState(patientId)

	fmt.Println("Patient Id is ", patientId, "and Patient details are ", string(bytes))

	if err != nil {
		return p, errors.New("Patient not found")
	}

	err = json.Unmarshal(bytes, &p)

	return p, nil

}

func checkIfAppointmentExists(stub shim.ChaincodeStubInterface, args []string) bool {

	patientBytes, err := stub.GetState(args[0])
	result := false

	if err == nil && patientBytes != nil {
		result = true
	}

	return result
}

func updateBalanceAPI(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	policyId := args[0]

	claimedAmountInt, _ := strconv.ParseInt(args[1], 10, 64)

	patientDetailsBytes, err := stub.GetState(policyId)

	if err != nil {
		return nil, errors.New("error while getting Policy Id")
	}

	patientDetails := Patient{}

	err = json.Unmarshal(patientDetailsBytes, &patientDetails)

	unavailedBalance := patientDetails.UnavailedBalance

	if unavailedBalance >= claimedAmountInt {
		unavailedBalance = unavailedBalance - claimedAmountInt
	} else {
		unavailedBalance = 0
	}

	patientDetails.UnavailedBalance = unavailedBalance
	patientDetails.ClaimedAmount = claimedAmountInt

	patBytes, _ := json.Marshal(patientDetails)

	err = stub.PutState(policyId, patBytes)

	return patBytes, nil

}

func createAppointment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	appointmentsExists := checkIfAppointmentExists(stub, args)
	fmt.Println("appointmentsExists : ", appointmentsExists)

	if !appointmentsExists {
		fmt.Println("Inside create appointment")

		patientDetails := Patient{}

		unavailedBalance, _ := strconv.ParseInt(args[7], 10, 64)
		claimedAmount, _ := strconv.ParseInt("0", 10, 64)
		contactNumber, _ := strconv.ParseInt(args[3], 10, 64)

		patientDetails.PloicyId = args[0]
		patientDetails.FirstName = args[1]
		patientDetails.LastName = args[2]
		patientDetails.Contact_Number = contactNumber
		patientDetails.Hospital = args[4]
		patientDetails.City = args[5]
		patientDetails.AppointmentTime = args[6]
		patientDetails.UnavailedBalance = unavailedBalance
		patientDetails.ClaimedAmount = claimedAmount

		patientDetailsBytes, err := json.Marshal(patientDetails)

		fmt.Println("Patient Details are : ", string(patientDetailsBytes))

		if err != nil {
			return nil, errors.New("Problem while saving Owner Details in BlockChain Network")

		}

		err = stub.PutState(patientDetails.PloicyId, patientDetailsBytes)

		//now owner has been added to block chain network, now we have to save the  Id as well

		bytes, err := stub.GetState("policy_Ids")

		var newPolicyId POLICY_ID_Holder

		err = json.Unmarshal(bytes, &newPolicyId)

		if err != nil {
			return nil, errors.New("error unmarshalling new Policy Id")
		}
		newPolicyId.POLICY_IDs = append(newPolicyId.POLICY_IDs, patientDetails.PloicyId)

		bytes, err = json.Marshal(newPolicyId)

		if err != nil {
			return nil, errors.New("error marshalling new Policy Id")
		}

		err = stub.PutState("policy_Ids", bytes)
		fmt.Println("Policy Id Saved is ", string(bytes))

		if err != nil {
			return nil, errors.New("Unable to put the state")
		}
	}
	return nil, nil
}
