package crypto

import(
    "fmt"
)

/*
@库描述：oi_tea算法的c语言版本移植

@说明：oi_tea原来的oi_symmetry_encrypt2加密时，某些填充字符是随机的，所以每次生成的加密串都不一样，
ex库增加了 oi_symmetry_encrypt2_regular 加密：
    1.填充字符固定，这样加密生成的字符串固定，可以用来生成 openid 等数据；
    2.regular版和普通版使用同一个解密函数
*/

/*
#ifndef WIN32
    typedef char BOOL;
#endif

#ifndef TRUE
#define TRUE 1
#endif

#ifndef FALSE
#define FALSE 0
#endif

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#ifdef WIN32
    #include <winsock2.h>

#define __bswap_constant_32(x) \
((((x) & 0xff000000u) >> 24) | (((x) & 0x00ff0000u) >>  8) |         \
(((x) & 0x0000ff00u) <<  8) | (((x) & 0x000000ffu) << 24))

# if __BYTE_ORDER == __BIG_ENDIAN
# define ntohl(x)    (x)
# define ntohs(x)    (x)
# define htonl(x)    (x)
# define htons(x)    (x)
# else
#  if __BYTE_ORDER == __LITTLE_ENDIAN
#   define ntohl(x)    __bswap_32 (x)
#   define ntohs(x)    __bswap_16 (x)
#   define htonl(x)    __bswap_32 (x)
#   define htons(x)    __bswap_16 (x)
#  endif
# endif

#else
    #include <netinet/in.h>
    #include <sys/time.h>
    #include <unistd.h>
#endif


void TeaEncryptECB(const char *pInBuf, const char *pKey, char *pOutBuf);
void TeaDecryptECB(const char *pInBuf, const char *pKey, char *pOutBuf);
void TeaEncryptECB3(const char *pInBuf, const char *pKey, char *pOutBuf);
void TeaDecryptECB3(const char *pInBuf, const char *pKey, char *pOutBuf);

void oi_symmetry_encrypt(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen);
BOOL oi_symmetry_decrypt(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen);

int oi_symmetry_encrypt2_len(int nInBufLen);
void oi_symmetry_encrypt2(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen);
void oi_symmetry_encrypt2_regular(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen);
BOOL oi_symmetry_decrypt2(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen);


typedef unsigned int WORD32;
const WORD32 DELTA = 0x9e3779b9;

#define ROUNDS 16
#define LOG_ROUNDS 4

void TeaEncryptECB(const char *pInBuf, const char *pKey, char *pOutBuf)
{
    WORD32 y, z;
    WORD32 sum;
    WORD32 k[4];
    int i;


    y = ntohl(*((WORD32*)pInBuf));
    z = ntohl(*((WORD32*)(pInBuf+4)));


    for ( i = 0; i<4; i++)
    {
        k[i] = ntohl(*((WORD32*)(pKey+i*4)));
    }

    sum = 0;
    for (i=0; i<ROUNDS; i++)
    {   
        sum += DELTA;
        y += ((z << 4) + k[0]) ^ (z + sum) ^ ((z >> 5) + k[1]);
        z += ((y << 4) + k[2]) ^ (y + sum) ^ ((y >> 5) + k[3]);
    }

    *((WORD32*)pOutBuf) = htonl(y);
    *((WORD32*)(pOutBuf+4)) = htonl(z);
}

void TeaDecryptECB(const char *pInBuf, const char *pKey, char *pOutBuf)
{
    WORD32 y, z, sum;
    WORD32 k[4];
    int i;

    y = ntohl(*((WORD32*)pInBuf));
    z = ntohl(*((WORD32*)(pInBuf+4)));

    for ( i=0; i<4; i++)
    {
        k[i] = ntohl(*((WORD32*)(pKey+i*4)));
    }

    sum = DELTA << LOG_ROUNDS;
    for (i=0; i<ROUNDS; i++)
    {
        z -= ((y << 4) + k[2]) ^ (y + sum) ^ ((y >> 5) + k[3]); 
        y -= ((z << 4) + k[0]) ^ (z + sum) ^ ((z >> 5) + k[1]);
        sum -= DELTA;
    }

    *((WORD32*)pOutBuf) = htonl(y);
    *((WORD32*)(pOutBuf+4)) = htonl(z);
}


void TeaEncryptECB3(const char *pInBuf, const char *pKey, char *pOutBuf)
{
    WORD32 y, z;
    WORD32 sum;
    WORD32 k[4];
    int i;


    y = ntohl(*((WORD32*)pInBuf));
    z = ntohl(*((WORD32*)(pInBuf+4)));

    for ( i = 0; i<4; i++)
    {
        k[i] = ntohl(*((WORD32*)(pKey+i*4)));
    }

    sum = 0;
    for (i=0; i<13; i++)
    {   
        sum += DELTA;
        y += ((z << 4) + k[0]) ^ (z + sum) ^ ((z >> 5) + k[1]);
        z += ((y << 4) + k[2]) ^ (y + sum) ^ ((y >> 5) + k[3]);
    }



    *((WORD32*)pOutBuf) = htonl(y);
    *((WORD32*)(pOutBuf+4)) = htonl(z);
}

void TeaDecryptECB3(const char *pInBuf, const char *pKey, char *pOutBuf)
{
    WORD32 y, z, sum;
    WORD32 k[4];
    int i;

    y = ntohl(*((WORD32*)pInBuf));
    z = ntohl(*((WORD32*)(pInBuf+4)));

    for ( i=0; i<4; i++)
    {
        k[i] = ntohl(*((WORD32*)(pKey+i*4)));
    }


    sum = DELTA << 3;
    for (i=1; i<=5; i++)
    {
        sum += DELTA;
    }


    for (i=0; i<13; i++)
    {
        z -= ((y << 4) + k[2]) ^ (y + sum) ^ ((y >> 5) + k[3]); 
        y -= ((z << 4) + k[0]) ^ (z + sum) ^ ((z >> 5) + k[1]);
        sum -= DELTA;
    }

    *((WORD32*)pOutBuf) = htonl(y);
    *((WORD32*)(pOutBuf+4)) = htonl(z);
}


#define SALT_LEN 2
#define ZERO_LEN 7

void oi_symmetry_encrypt(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen)
{
    int nPadSaltBodyZeroLen;
    int nPadlen;
    char src_buf[8], zero_iv[8], *iv_buf;
    int src_i, i, j;

    nPadSaltBodyZeroLen = nInBufLen+1+SALT_LEN+ZERO_LEN;
    if((nPadlen=(nPadSaltBodyZeroLen%8)))
    {
        nPadlen=8-nPadlen;
    }

    src_buf[0] = (((char)rand()) & 0x0f8) | (char)nPadlen;
    src_i = 1;

    while(nPadlen--)
        src_buf[src_i++]=(char)rand(); 


    memset(zero_iv, 0, 8);
    iv_buf = zero_iv; 

    *pOutBufLen = 0; 

    for (i=1;i<=SALT_LEN;)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=(char)rand();
            i++; 
        }

        if (src_i==8)
        {

            
            for (j=0;j<8;j++)
                src_buf[j]^=iv_buf[j];

            TeaEncryptECB(src_buf, pKey, pOutBuf);
            src_i=0;
            iv_buf=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }

    while(nInBufLen)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=*(pInBuf++);
            nInBufLen--;
        }

        if (src_i==8)
        {

            for (i=0;i<8;i++)
                src_buf[i]^=iv_buf[i];
            
            TeaEncryptECB(src_buf, pKey, pOutBuf);
            src_i=0;
            iv_buf=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }

    for (i=1;i<=ZERO_LEN;)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=0;
            i++;
        }

        if (src_i==8)
        {

            
            for (j=0;j<8;j++) 
                src_buf[j]^=iv_buf[j];
            
            TeaEncryptECB(src_buf, pKey, pOutBuf);
            src_i=0;
            iv_buf=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }
}


BOOL oi_symmetry_decrypt(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen)
{
    int nPadLen, nPlainLen;
    char dest_buf[8];
    const char *iv_buf;
    int dest_i, i, j;


    if ((nInBufLen%8) || (nInBufLen<16)) return FALSE;


    TeaDecryptECB(pInBuf, pKey, dest_buf);

    nPadLen = dest_buf[0] & 0x7;


    i = nInBufLen-1-nPadLen-SALT_LEN-ZERO_LEN;
    if (*pOutBufLen<i) return FALSE;
    *pOutBufLen = i;
    if (*pOutBufLen < 0) return FALSE;


    iv_buf = pInBuf;
    nInBufLen -= 8;
    pInBuf += 8;

    dest_i=1;

    dest_i+=nPadLen;

    for (i=1; i<=SALT_LEN;)
    {
        if (dest_i<8)
        {
            dest_i++;
            i++;
        }

        if (dest_i==8)
        {
            TeaDecryptECB(pInBuf, pKey, dest_buf);
            for (j=0; j<8; j++)
                dest_buf[j]^=iv_buf[j];
        
            iv_buf = pInBuf;
            nInBufLen -= 8;
            pInBuf += 8;

            dest_i=0;
        }
    }

    nPlainLen=*pOutBufLen;
    while (nPlainLen)
    {
        if (dest_i<8)
        {
            *(pOutBuf++)=dest_buf[dest_i++];
            nPlainLen--;
        }
        else if (dest_i==8)
        {
            TeaDecryptECB(pInBuf, pKey, dest_buf);
            for (i=0; i<8; i++)
                dest_buf[i]^=iv_buf[i];
        
            iv_buf = pInBuf;
            nInBufLen -= 8;
            pInBuf += 8;

            dest_i=0; 
        }
    }

    for (i=1;i<=ZERO_LEN;)
    {
        if (dest_i<8)
        {
            if(dest_buf[dest_i++]) return FALSE;
            i++;
        }
        else if (dest_i==8)
        {
            TeaDecryptECB(pInBuf, pKey, dest_buf);
            for (j=0; j<8; j++)
                dest_buf[j]^=iv_buf[j];
        
            iv_buf = pInBuf;
            nInBufLen -= 8;
            pInBuf += 8;

            dest_i=0; 
        }

    }

    return TRUE;
}

int oi_symmetry_encrypt2_len(int nInBufLen)
{
    int nPadSaltBodyZeroLen;
    int nPadlen;

    nPadSaltBodyZeroLen = nInBufLen+1+SALT_LEN+ZERO_LEN;
    if((nPadlen=(nPadSaltBodyZeroLen%8))) 
    {

        nPadlen=8-nPadlen;
    }

    return nPadSaltBodyZeroLen+nPadlen;
}


void oi_symmetry_encrypt2(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen)
{
    int nPadSaltBodyZeroLen;
    int nPadlen;
    char src_buf[8], iv_plain[8], *iv_crypt;
    int src_i, i, j;
    
    static int iSrand = 0;
    if (iSrand == 0)
    {
        srand(time(NULL));// * getpid());
        iSrand = 1;
    }

    nPadSaltBodyZeroLen = nInBufLen+1+SALT_LEN+ZERO_LEN;
    if((nPadlen=(nPadSaltBodyZeroLen%8)))
    {
        nPadlen=8-nPadlen;
    }

    src_buf[0] = (((char)rand()) & 0x0f8) | (char)nPadlen;
    src_i = 1;

    while(nPadlen--)
        src_buf[src_i++]=(char)rand();

    for ( i=0; i<8; i++)
        iv_plain[i] = 0;
    iv_crypt = iv_plain;

    *pOutBufLen = 0;

    for (i=1;i<=SALT_LEN;)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=(char)rand();
            i++; 
        }

        if (src_i==8)
        {
            for (j=0;j<8;j++) 
                src_buf[j]^=iv_crypt[j];

            TeaEncryptECB(src_buf, pKey, pOutBuf);

            for (j=0;j<8;j++)
                pOutBuf[j]^=iv_plain[j];

            for (j=0;j<8;j++)
                iv_plain[j]=src_buf[j];

            src_i=0;
            iv_crypt=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }

    while(nInBufLen)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=*(pInBuf++);
            nInBufLen--;
        }

        if (src_i==8)
        {
            for (j=0;j<8;j++) 
                src_buf[j]^=iv_crypt[j];

            TeaEncryptECB(src_buf, pKey, pOutBuf);

            for (j=0;j<8;j++) 
                pOutBuf[j]^=iv_plain[j];

            for (j=0;j<8;j++)
                iv_plain[j]=src_buf[j];

            src_i=0;
            iv_crypt=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }

    for (i=1;i<=ZERO_LEN;)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=0;
            i++;
        }

        if (src_i==8)
        {
            for (j=0;j<8;j++) 
                src_buf[j]^=iv_crypt[j];

            TeaEncryptECB(src_buf, pKey, pOutBuf);

            for (j=0;j<8;j++)
                pOutBuf[j]^=iv_plain[j];

            for (j=0;j<8;j++)
                iv_plain[j]=src_buf[j];

            src_i=0;
            iv_crypt=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }
}


void oi_symmetry_encrypt2_regular(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen)
{
    int nPadSaltBodyZeroLen;
    int nPadlen;
    char src_buf[8], iv_plain[8], *iv_crypt;
    int src_i, i, j;
    char rand_buf = '1';

    nPadSaltBodyZeroLen = nInBufLen+1+SALT_LEN+ZERO_LEN;
    if((nPadlen=nPadSaltBodyZeroLen%8))
    {
        nPadlen=8-nPadlen;
    }

    src_buf[0] = (rand_buf & 0x0f8) | (char)nPadlen;
    src_i = 1;

    while(nPadlen--)
        src_buf[src_i++]=rand_buf;

    for ( i=0; i<8; i++)
        iv_plain[i] = 0;
    iv_crypt = iv_plain;

    *pOutBufLen = 0;

    for (i=1;i<=SALT_LEN;)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=rand_buf;
            i++;
        }

        if (src_i==8)
        {
            for (j=0;j<8;j++)
                src_buf[j]^=iv_crypt[j];

            TeaEncryptECB(src_buf, pKey, pOutBuf);

            for (j=0;j<8;j++)
                pOutBuf[j]^=iv_plain[j];

            for (j=0;j<8;j++)
                iv_plain[j]=src_buf[j];

            src_i=0;
            iv_crypt=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }

    while(nInBufLen)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=*(pInBuf++);
            nInBufLen--;
        }

        if (src_i==8)
        {
            for (j=0;j<8;j++)
                src_buf[j]^=iv_crypt[j];

            TeaEncryptECB(src_buf, pKey, pOutBuf);

            for (j=0;j<8;j++)
                pOutBuf[j]^=iv_plain[j];

            for (j=0;j<8;j++)
                iv_plain[j]=src_buf[j];

            src_i=0;
            iv_crypt=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }

    for (i=1;i<=ZERO_LEN;)
    {
        if (src_i<8)
        {
            src_buf[src_i++]=0;
            i++;
        }

        if (src_i==8)
        {
            for (j=0;j<8;j++)
                src_buf[j]^=iv_crypt[j];

            TeaEncryptECB(src_buf, pKey, pOutBuf);

            for (j=0;j<8;j++)
                pOutBuf[j]^=iv_plain[j];

            for (j=0;j<8;j++)
                iv_plain[j]=src_buf[j];

            src_i=0;
            iv_crypt=pOutBuf;
            *pOutBufLen+=8;
            pOutBuf+=8;
        }
    }
}


BOOL oi_symmetry_decrypt2(const char* pInBuf, int nInBufLen, const char* pKey, char* pOutBuf, int *pOutBufLen)
{
    int nPadLen, nPlainLen;
    char dest_buf[8], zero_buf[8];
    const char *iv_pre_crypt, *iv_cur_crypt;
    int dest_i, i, j;

    int nBufPos;
    nBufPos = 0;

    if ((nInBufLen%8) || (nInBufLen<16)) return FALSE;

    TeaDecryptECB(pInBuf, pKey, dest_buf);

    nPadLen = dest_buf[0] & 0x7;

    i = nInBufLen-1-nPadLen-SALT_LEN-ZERO_LEN;
    if ((*pOutBufLen<i) || (i<0)) return FALSE;
    *pOutBufLen = i;

    for ( i=0; i<8; i++)
        zero_buf[i] = 0;

    iv_pre_crypt = zero_buf;
    iv_cur_crypt = pInBuf;

    pInBuf += 8;
    nBufPos += 8;

    dest_i=1;

    dest_i+=nPadLen;

    for (i=1; i<=SALT_LEN;)
    {
        if (dest_i<8)
        {
            dest_i++;
            i++;
        }
        else if (dest_i==8)
        {
            iv_pre_crypt = iv_cur_crypt;
            iv_cur_crypt = pInBuf; 

            for (j=0; j<8; j++)
            {
                if( (nBufPos + j) >= nInBufLen)
                    return FALSE;
                dest_buf[j]^=pInBuf[j];
            }

            TeaDecryptECB(dest_buf, pKey, dest_buf);

            pInBuf += 8;
            nBufPos += 8;

            dest_i=0;
        }
    }

    nPlainLen=*pOutBufLen;
    while (nPlainLen)
    {
        if (dest_i<8)
        {
            *(pOutBuf++)=dest_buf[dest_i]^iv_pre_crypt[dest_i];
            dest_i++;
            nPlainLen--;
        }
        else if (dest_i==8)
        {
            iv_pre_crypt = iv_cur_crypt;
            iv_cur_crypt = pInBuf; 

            for (j=0; j<8; j++)
            {
                if( (nBufPos + j) >= nInBufLen)
                    return FALSE;
                dest_buf[j]^=pInBuf[j];
            }

            TeaDecryptECB(dest_buf, pKey, dest_buf);

            pInBuf += 8;
            nBufPos += 8;

            dest_i=0;
        }
    }

    for (i=1;i<=ZERO_LEN;)
    {
        if (dest_i<8)
        {
            if(dest_buf[dest_i]^iv_pre_crypt[dest_i]) return FALSE;
            dest_i++;
            i++;
        }
        else if (dest_i==8)
        {
            iv_pre_crypt = iv_cur_crypt;
            iv_cur_crypt = pInBuf; 

            for (j=0; j<8; j++)
            {
                if( (nBufPos + j) >= nInBufLen)
                    return FALSE;
                dest_buf[j]^=pInBuf[j];
            }

            TeaDecryptECB(dest_buf, pKey, dest_buf);
            
            pInBuf += 8;
            nBufPos += 8;
            dest_i=0;
        }

    }

    return TRUE;
}
*/
import "C"

/*pOutBuffer、pInBuffer均为8byte, pKey为16byte*/
func TeaEncryptECB(inbuf []byte, key []byte) []byte {
    var outbuf [16]C.char

    C.TeaEncryptECB(C.CString(string(inbuf)), C.CString(string(key)), &outbuf[0])

    out := C.GoString(&outbuf[0])
    return []byte(out)
}

/*pOutBuffer、pInBuffer均为8byte, pKey为16byte*/
func TeaDecryptECB(inbuf []byte, key []byte) []byte {
    var outbuf [16]C.char

    C.TeaDecryptECB(C.CString(string(inbuf)), C.CString(string(key)), &outbuf[0])

    out := C.GoString(&outbuf[0])
    return []byte(out)
}

/*pOutBuffer、pInBuffer均为8byte, pKey为16byte*/
func TeaEncryptECB3(inbuf []byte, key []byte) []byte {
    var outbuf [16]C.char

    C.TeaEncryptECB3(C.CString(string(inbuf)), C.CString(string(key)), &outbuf[0])

    out := C.GoString(&outbuf[0])
    return []byte(out)
}

/*pOutBuffer、pInBuffer均为8byte, pKey为16byte*/
func TeaDecryptECB3(inbuf []byte , key []byte) []byte {
    var outbuf [16]C.char

    C.TeaDecryptECB3(C.CString(string(inbuf)), C.CString(string(key)), &outbuf[0])

    out := C.GoString(&outbuf[0])
    return []byte(out)
}

const(
    MAX_OUT_BUF_LEN = 10240
    TEA_ENC_KEY_LEN = 16
)

const(
    FALSE = 0
    TRUE = 1
)

/*pKey为16byte*/
/*
	输入:pInBuf为需加密的明文部分(Body),nInBufLen为pInBuf长度;
	输出:pOutBuf为密文格式,pOutBufLen为pOutBuf的长度是8byte的倍数;
*/
/*TEA加密算法,CBC模式*/
/*密文格式:PadLen(1byte)+Padding(var,0-7byte)+Salt(2byte)+Body(var byte)+Zero(7byte)*/
func Oi_symmetry_encrypt(inbuf []byte, key []byte) ([]byte, error) {
    if len(key) != TEA_ENC_KEY_LEN {
        return nil, fmt.Errorf("invalid key(\"%s\"), length:%d", string(key), len(key))
    }

    inbuflen := len(inbuf)
    var outbuf [MAX_OUT_BUF_LEN]C.char
    var outbuflen C.int

    C.oi_symmetry_encrypt(C.CString(string(inbuf)), C.int(inbuflen), C.CString(string(key)), &outbuf[0], &outbuflen)

    out := C.GoStringN(&outbuf[0], outbuflen)

    return []byte(out)[:int(outbuflen)], nil
}

/*pKey为16byte*/
/*
	输入:pInBuf为密文格式,nInBufLen为pInBuf的长度是8byte的倍数;
	输出:pOutBuf为明文(Body),pOutBufLen为pOutBuf的长度;
	返回值:如果格式正确返回TRUE;
*/
/*TEA解密算法,CBC模式*/
/*密文格式:PadLen(1byte)+Padding(var,0-7byte)+Salt(2byte)+Body(var byte)+Zero(7byte)*/
func Oi_symmetry_decrypt(inbuf []byte, key []byte) ([]byte, error) {
    if len(key) != TEA_ENC_KEY_LEN {
        return nil, fmt.Errorf("invalid key(\"%s\"), length:%d", string(key), len(key))
    }

    inbuflen := len(inbuf)
    var outbuf [MAX_OUT_BUF_LEN]C.char
    var outbuflen C.int

    var ret C.BOOL = C.oi_symmetry_decrypt(C.CString(string(inbuf)), C.int(inbuflen), C.CString(string(key)), &outbuf[0], &outbuflen)
    if int(ret) != TRUE {
        return nil, fmt.Errorf("decrypt fail:%d", int(ret))
    }
    
    out := C.GoStringN(&outbuf[0], outbuflen)

    return []byte(out)[:int(outbuflen)], nil
}

/*pKey为16byte*/
/*
	输入:nInBufLen为需加密的明文部分(Body)长度;
	输出:返回为加密后的长度(是8byte的倍数);
*/
/*TEA加密算法,CBC模式*/
/*密文格式:PadLen(1byte)+Padding(var,0-7byte)+Salt(2byte)+Body(var byte)+Zero(7byte)*/
func Oi_symmetry_encrypt2_len(inbuflen int) int {
    ret := C.oi_symmetry_encrypt2_len(C.int(inbuflen))
    return int(ret)
}


/*pKey为16byte*/
/*
	输入:pInBuf为需加密的明文部分(Body),nInBufLen为pInBuf长度;
	输出:pOutBuf为密文格式,pOutBufLen为pOutBuf的长度是8byte的倍数;
*/
/*TEA加密算法,CBC模式*/
/*密文格式:PadLen(1byte)+Padding(var,0-7byte)+Salt(2byte)+Body(var byte)+Zero(7byte)*/
func Oi_symmetry_encrypt2(inbuf []byte, key []byte) ([]byte, error) {
    if len(key) != TEA_ENC_KEY_LEN {
        return nil, fmt.Errorf("invalid key(\"%s\"), length:%d", string(key), len(key))
    }

    inbuflen := len(inbuf)
    outLen := Oi_symmetry_encrypt2_len(inbuflen)
    if outLen > MAX_OUT_BUF_LEN {
        return nil, fmt.Errorf("out buf len too large:%d > %d", outLen, MAX_OUT_BUF_LEN)
    }

    var outbuf [MAX_OUT_BUF_LEN]C.char
    var outbuflen C.int

    C.oi_symmetry_encrypt2(C.CString(string(inbuf)), C.int(inbuflen), C.CString(string(key)), &outbuf[0], &outbuflen)

    out := C.GoStringN(&outbuf[0], outbuflen)

    return []byte(out)[:int(outbuflen)], nil
}

// 新增函数，填充字符不再随机，导致生成的字符串固定
func Oi_symmetry_encrypt2_regular(inbuf []byte, key []byte) ([]byte, error) {
    if len(key) != TEA_ENC_KEY_LEN {
        return nil, fmt.Errorf("invalid key(\"%s\"), length:%d", string(key), len(key))
    }

    inbuflen := len(inbuf)
    outLen := Oi_symmetry_encrypt2_len(inbuflen)
    if outLen > MAX_OUT_BUF_LEN {
        return nil, fmt.Errorf("out buf len too large:%d > %d", outLen, MAX_OUT_BUF_LEN)
    }

    var outbuf [MAX_OUT_BUF_LEN]C.char
    var outbuflen C.int

    C.oi_symmetry_encrypt2_regular(C.CString(string(inbuf)), C.int(inbuflen), C.CString(string(key)), &outbuf[0], &outbuflen)

    out := C.GoStringN(&outbuf[0], outbuflen)

    return []byte(out)[:int(outbuflen)], nil
}


/*pKey为16byte*/
/*
	输入:pInBuf为密文格式,nInBufLen为pInBuf的长度是8byte的倍数; *pOutBufLen为接收缓冲区的长度
		特别注意*pOutBufLen应预置接收缓冲区的长度!
	输出:pOutBuf为明文(Body),pOutBufLen为pOutBuf的长度,至少应预留nInBufLen-10;
	返回值:如果格式正确返回TRUE;
*/
/*TEA解密算法,CBC模式*/
/*密文格式:PadLen(1byte)+Padding(var,0-7byte)+Salt(2byte)+Body(var byte)+Zero(7byte)*/
func Oi_symmetry_decrypt2(inbuf []byte, key []byte) ([]byte, error) {
    if len(key) != TEA_ENC_KEY_LEN {
        return nil, fmt.Errorf("invalid key(\"%s\"), length:%d", string(key), len(key))
    }

    inbuflen := len(inbuf)
    var outbuf [MAX_OUT_BUF_LEN]C.char
    outbuflen := C.int(int(inbuflen-10))

    var ret C.BOOL = C.oi_symmetry_decrypt2(C.CString(string(inbuf)), C.int(inbuflen), C.CString(string(key)), &outbuf[0], &outbuflen)
    if int(ret) != TRUE {
        return nil, fmt.Errorf("decrypt fail:%d", int(ret))
    }
    
    out := C.GoStringN(&outbuf[0], outbuflen)

    return []byte(out)[:int(outbuflen)], nil
}
