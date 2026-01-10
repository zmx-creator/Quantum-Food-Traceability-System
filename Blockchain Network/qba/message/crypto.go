package message

import (
	//"crypto/sha256"
	//"encoding/hex"
	"encoding/json"
	"log"
	//"strings"
)

const(
    KeyLength = 64    //key length
    MessageLength = 4096   //maeesage length
)

func Digest(obj interface{}) ([]byte, error) {
	
	result := make([]byte, MessageLength)

	
	content, err := json.Marshal(obj)
	if err != nil {
		log.Printf("[Encode] marshal error: %v", err)
		return nil, err
	}


	for i := 0; i < MessageLength; i++ {
		if i < len(content)*8 { 
			byteIdx := i / 8
			bitPos := uint(7 - (i % 8))
			result[i] = (content[byteIdx] >> bitPos) & 1
		} else {
			
			result[i] = 0
		}
	}

	return result, nil
}



func Sign(message []byte, polynomial []byte,signerKeys_X []byte,signerKeys_Y []byte,signerKeys_Z []byte) QDSSignature {   
    //generate QDS signature
    //use user's QDS key x,y,z and a irreducible polynomial to generate a signature
    //signature.Digest:signature content is byte array composed of binary digits
    //signature.Factorsignature also take irreducible polynomial factor
        
	
    //generate hash matrix
    matrix := generateToeplitzMatrix(signerKeys_X, polynomial)
	
    //compute hash value
    hashValue := computeHash(matrix, message)

    //log output the hash value calculated by signer  
    log.Printf("Signer Calculate Hash Value: %v", hashValue)

    //generate signature
    dig := make([]byte, KeyLength)
    p := make([]byte, KeyLength)

    for i := 0; i < KeyLength; i++ {
        dig[i] = (hashValue[i] + signerKeys_Y[i]) % 2
        p[i] = (polynomial[i] + signerKeys_Z[i]) % 2
    }

    signature := QDSSignature{
        Digest: dig,
        Factor:   p,
    }

    log.Printf("Generate Sign: Digest:%v, Factor:%v", signature.Digest, signature.Factor)

    return signature
}

func generateToeplitzMatrix(X []byte, polynomial []byte) [][]byte {
    // generate Toeplitz Hash Matrix
    matrix := make([][]byte, KeyLength)
    for i := range matrix {
        matrix[i] = make([]byte, MessageLength)
    }
    //first coloum is X
    for i := 0; i < KeyLength; i++ {
        matrix[i][0] = X[i]
    }
    
    //follow coloum generate by polynomial transformation
    for j := 1; j < MessageLength; j++ {
        //calculate the dot product of the polynomial and the previous column
        temp := 0
        for k := 0; k < KeyLength; k++ {
            temp += int(polynomial[k] * matrix[k][j-1])
        }
        temp = temp % 2
        
        //generate new column
        matrix[0][j] = byte(temp)
        for i := 1; i < 8; i++ {
            matrix[i][j] = matrix[i-1][j-1]
        }
    }

    return matrix
}


//compute Hash Value
func computeHash(matrix [][]byte, message []byte) []byte {
    hash := make([]byte, KeyLength)

    for i := 0; i < KeyLength; i++ {
        sum := 0
        for j := 0; j < MessageLength; j++ {
            sum += int(matrix[i][j] * message[j])
        }
        hash[i] = byte(sum % 2)
    }

    return hash
}


//verify QDS key
func QDSVerify(message []byte, signature QDSSignature, forwarderKeys_X []byte ,forwarderKeys_Y []byte ,forwarderKeys_Z []byte , verifierKeys_X []byte,verifierKeys_Y []byte,verifierKeys_Z []byte) string {
	// calculate to restore key
	Xr := make([]byte, KeyLength)
	Yr := make([]byte, KeyLength)
	Zr := make([]byte, KeyLength)

	for i := 0; i < KeyLength; i++ {
		Xr[i] = (forwarderKeys_X[i] + verifierKeys_X[i]) % 2
		Yr[i] = (forwarderKeys_Y[i] + verifierKeys_Y[i]) % 2
		Zr[i] = (forwarderKeys_Z[i] + verifierKeys_Z[i]) % 2
	}

	log.Printf("restore key Xr:%v, Yr:%v, Zr:%v", Xr, Yr, Zr)

	// decode signatore
	decodedHash := make([]byte, KeyLength)
	decodedPolynomial := make([]byte, KeyLength)

	for i := 0; i < KeyLength; i++ {
		decodedHash[i] = (signature.Digest[i] + Yr[i]) % 2
		decodedPolynomial[i] = (signature.Factor[i] + Zr[i]) % 2
	}

	log.Printf("decode to get hash value: %v", decodedHash)
	log.Printf("decode to get polynomial: %v", decodedPolynomial)

	//re-compute hash matrix
	matrix := generateToeplitzMatrix(Xr, decodedPolynomial)

	// calculate new hash value
	newHash := computeHash(matrix, message)

	log.Printf("hash value calculated by verifier: %v", newHash)

	// campare hash value
	for i := 0; i < KeyLength; i++ {
		if decodedHash[i] != newHash[i] {
			log.Printf("verify failed:hash value doesn't match")
			return "verify failed"
		}
	}

	log.Printf("verify succeed:hash value match")
	return "verify succeed"
}
